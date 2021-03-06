package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

func parseListing(listing *Listing) {

	// default is no username or password, because the request may be invalid
	listing.username = ""
	listing.password = ""

	// debug base request xml
	fmt.Println(string(listing.Request))

	// validate request
	// this means checking it has a username, password and unique ID
	var baseRequest requestBase
	if err := xml.Unmarshal(listing.Request, &baseRequest); err != nil {
		log.Println("invalid request format")
		listing.Response, _ = xml.Marshal(ResponseErrorInvalid)
		return
	}

	// failed auth
	if !checkAuth(baseRequest.Username, baseRequest.Password) {
		log.Printf("invalid auth for request %d\n", baseRequest.ID)
		listing.username = ""
		listing.password = ""
		listing.Response, _ = xml.Marshal(ResponseErrorAuth)
		return
	}

	// check for duplicate ID
	for _, slot := range usersIDs[baseRequest.Username] {
		if baseRequest.ID == slot {
			log.Printf("duplicate id %d\n", baseRequest.ID)
			listing.Response, _ = xml.Marshal(ResponseErrorDuplicateID)
			return
		}
	}

	// save ID
	usersIDs[baseRequest.Username] = append(
		usersIDs[baseRequest.Username], baseRequest.ID)

	// debug base request
	log.Printf("processing type %s, id %d for %s\n",
		baseRequest.XMLName.Local,
		baseRequest.ID,
		baseRequest.Username)

	// if format is valid, assign username and password to listing
	// they will later be checked against the server users for the
	// proper request types
	listing.username = baseRequest.Username
	listing.password = baseRequest.Password

	// validate and handle request
	switch baseRequest.XMLName.Local {
	// --------------------------------------------------------------------
	case "availability":
		var req RequestAvailability

		// fully parse request
		if err := xml.Unmarshal(listing.Request, &req); err != nil {
			log.Println("invalid request format")
			listing.username = ""
			listing.password = ""
			listing.Response, _ = xml.Marshal(ResponseErrorInvalid)
			return
		}

		// prepare response
		var res ResponseAvailability
		res.XMLName = xml.Name{Local: "response"}
		res.Code = 200

		// get all free slots
		for idx := range slots {
			if slots[idx] == "free" {
				res.Slots = append(res.Slots, idx+1)
			}
		}

		// create response payload
		resBinary, err := xml.Marshal(res)
		if err != nil {
			panic(err)
		}

		// assign response to listing
		listing.Response = resBinary

		log.Printf("parsed request %d of type '%s'\n",
			req.ID, req.XMLName.Local)

		return
	// --------------------------------------------------------------------
	case "reserve":
		var req RequestReserve

		// fully parse request
		if err := xml.Unmarshal(listing.Request, &req); err != nil {
			log.Println("invalid request format")
			listing.username = ""
			listing.password = ""
			listing.Response, _ = xml.Marshal(ResponseErrorInvalid)
			return
		}

		// check too many booked slots
		if len(usersSlots[req.Username]) >= maxBookedSlots {
			log.Printf("max booked slots, can't reserve more %d/%d\n",
				len(usersSlots[req.Username]), maxBookedSlots)

			var resErr ResponseError

			resErr.XMLName = xml.Name{Local: "response"}
			resErr.Code = 409 // limit reached
			resErr.Body = fmt.Sprintf(
				"Reservation failed, you already hold the maximum permitted number of reservations - %d",
				maxBookedSlots)
			listing.Response, _ = xml.Marshal(resErr)

			return
		}

		// check if invalid slot
		if req.SlotID < 1 || req.SlotID > len(slots) {
			log.Printf("invalid slot %d\n", req.SlotID)

			var resErr ResponseError

			resErr.XMLName = xml.Name{Local: "response"}
			resErr.Code = 403 // slot does not exist
			resErr.Body = fmt.Sprintf("Slot %d does not exist", req.SlotID)
			listing.Response, _ = xml.Marshal(resErr)

			return
		}

		// check slot already taken
		if slots[req.SlotID-1] != "free" {
			log.Printf("invalid slot %d, not free\n", req.SlotID)

			var resErr ResponseError

			resErr.XMLName = xml.Name{Local: "response"}
			resErr.Code = 409 // also slot is not free
			resErr.Body = fmt.Sprintf("Slot %d is not free.", req.SlotID)
			listing.Response, _ = xml.Marshal(resErr)

			return
		}

		// assign slot
		slots[req.SlotID-1] = req.Username
		usersSlots[req.Username] = append(usersSlots[req.Username], req.SlotID)

		// all OK, prepare response
		var res ResponseReserve
		res.XMLName = xml.Name{Local: "response"}
		res.Code = 200
		res.Reserve = fmt.Sprintf("Reserved slot %d", req.SlotID)

		// create response payload
		resBinary, err := xml.Marshal(res)
		if err != nil {
			panic(err)
		}

		// assign response to listing
		listing.Response = resBinary

		log.Printf("parsed request %d of type '%s'\n",
			req.ID, req.XMLName.Local)

		return
	// --------------------------------------------------------------------
	case "cancel":
		var req RequestCancel

		// fully parse request
		if err := xml.Unmarshal(listing.Request, &req); err != nil {
			log.Println("invalid request format")
			listing.username = ""
			listing.password = ""
			listing.Response, _ = xml.Marshal(ResponseErrorInvalid)
			return
		}

		// check if invalid slot
		// this is not specified in the docs, but it makes sense right?
		// you shouldn't be able to cancel a slot that does not exist
		if req.SlotID < 1 || req.SlotID > len(slots) {
			log.Printf("invalid slot %d\n", req.SlotID)

			var resErr ResponseError

			resErr.XMLName = xml.Name{Local: "response"}
			resErr.Code = 403 // slot does not exist
			resErr.Body = fmt.Sprintf("Slot %d does not exist", req.SlotID)
			listing.Response, _ = xml.Marshal(resErr)

			return
		}

		// check slot is owned by user
		// note that the error message is changed from the one specified in the
		// labscript, but the code is the same
		// the message is just for the human to read, your agent should process
		// the error code
		if slots[req.SlotID-1] != req.Username {
			log.Printf("invalid slot %d, not yours\n", req.SlotID)

			var resErr ResponseError

			resErr.XMLName = xml.Name{Local: "response"}
			resErr.Code = 409 // also slot is not free
			resErr.Body = fmt.Sprintf(
				"Cancel failed, slot %d was not reserved by you",
				req.SlotID)
			listing.Response, _ = xml.Marshal(resErr)

			return
		}

		// mark slot as free
		slots[req.SlotID-1] = "free"
		// remove item from user slots
		for idx, slot := range usersSlots[req.Username] {
			if slot == req.SlotID {
				usersSlots[req.Username] = remove(usersSlots[req.Username], idx)
			}
		}

		// all OK, prepare response
		var res ResponseBasic
		res.XMLName = xml.Name{Local: "response"}
		res.Code = 200
		res.Body = fmt.Sprintf(
			"The reservation for slot %d has been cancelled",
			req.SlotID)

		// create response payload
		resBinary, err := xml.Marshal(res)
		if err != nil {
			panic(err)
		}

		// assign response to listing
		listing.Response = resBinary

		log.Printf("parsed request %d of type '%s'\n",
			req.ID, req.XMLName.Local)

		return
	// --------------------------------------------------------------------
	case "bookings":
		var req RequestBookings

		// fully parse request
		if err := xml.Unmarshal(listing.Request, &req); err != nil {
			log.Println("invalid request format")
			listing.username = ""
			listing.password = ""
			listing.Response, _ = xml.Marshal(ResponseErrorInvalid)
			return
		}

		// prepare response
		var res ResponseBookings
		res.XMLName = xml.Name{Local: "response"}
		res.Code = 200

		// get all booked slots for this user
		for idx := range usersSlots[req.Username] {
			res.Slots = append(res.Slots, idx+1)
		}

		// create response payload
		resBinary, err := xml.Marshal(res)
		if err != nil {
			panic(err)
		}

		// assign response to listing
		listing.Response = resBinary

		log.Printf("parsed request %d of type '%s'\n",
			req.ID, req.XMLName.Local)

		return
	// --------------------------------------------------------------------
	default:
		log.Fatalln("should not reach this")
	}
}
