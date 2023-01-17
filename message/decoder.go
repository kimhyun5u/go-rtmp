package message

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/yutopp/amf0-go"
)

type Decoder struct {
	r      io.Reader
	typeID byte
}

func NewDecoder(r io.Reader, typeID byte) *Decoder {
	return &Decoder{
		r:      r,
		typeID: typeID,
	}
}

func (dec *Decoder) Decode(msg *Message) error {
	switch TypeID(dec.typeID) {
	case TypeIDAudioMessage:
		return dec.decodeAudioMessage(msg)
	case TypeIDVideoMessage:
		return dec.decodeVideoMessage(msg)
	case TypeIDDataMessageAMF0:
		return dec.decodeDataMessage(msg)
	case TypeIDCommandMessageAMF0:
		return dec.decodeCommandMessage(msg)
	default:
		return fmt.Errorf("unexpected message type: %d", dec.typeID)

	}
}

func (dec *Decoder) decodeAudioMessage(msg *Message) error {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, dec.r)
	if err != nil {
		return err
	}

	*msg = &AudioMessage{
		Payload: buf.Bytes(),
	}

	return nil
}

func (dec *Decoder) decodeVideoMessage(msg *Message) error {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, dec.r)
	if err != nil {
		return err
	}

	*msg = &VideoMessage{
		Payload: buf.Bytes(),
	}

	return nil
}

func (dec *Decoder) decodeDataMessage(msg *Message) error {
	d := amf0.NewDecoder(dec.r)

	var name string
	if err := d.Decode(&name); err != nil {
		return err
	}
	log.Printf("name = %+v", name)

	var data interface{}
	switch name {
	case "onMetaData":
		var metadata map[string]interface{}
		if err := d.Decode(&metadata); err != nil {
			return err
		}
		log.Printf("onMetaData: metadata = %+v", metadata)
		data = &NetStreamOnMetaData{
			MetaData: metadata,
		}
	case "@setDataFram":
		log.Println("Ignored data message: @setDataFrame")
	default:
		return errors.New("Not supported data message: " + name)
	}

	*msg = &DataMEssageAMF0{
		Name: name,
		Data: data,
	}

	return nil
}

func (dec *Decoder) decodeCommandMessage(msg *Message) error {
	d := amf0.NewDecoder(dec.r)

	var name string
	if err := d.Decode(&name); err != nil {
		return err
	}
	log.Printf("name = %+v", name)

	var transactionID int64
	if err := d.Decode(&transactionID); err != nil {
		return err
	}

	log.Printf("transactionID = %+v", transactionID)

	var command interface{}
	switch name {
	case "connect":
		var object map[string]interface{}
		if err := d.Decode(&object); err != nil {
			return err
		}
		log.Printf("command: object = %+v", object)
		command = &NetConnectionConnection{
			CommandObject: object,
		}
	case "releaseStream":
		log.Printf("ignored releaseStream")
	case "FCPublish":
		log.Printf("ignored FCPublish")
	case "createStream":
		var object interface{}
		if err := d.Decode(&object); err != nil {
			return err
		}
		if object == nil {
			break
		}
		command = &NetConnectionCreateStream{}
	case "publish":
		var commandObject interface{}
		if err := d.Decode(&commandObject); err != nil {
			return err
		}
		var publishingName string
		if err := d.Decode(&publishingName); err != nil {
			return err
		}
		var publishingType string
		if err := d.Decode(&publishingType); err != nil {
			return err
		}
		command = &NetStreamPublish{
			CommandObject:  commandObject,
			PublishingName: publishingName,
			PublishingType: publishingType,
		}
	default:
		return errors.New("Not supported command: " + name)
	}

	*msg = &CommandMessageAMF0{
		CommandName:   name,
		TransactionID: transactionID,
		Command:       command,
	}

	return nil
}
