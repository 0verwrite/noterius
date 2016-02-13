package land

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Nyarum/noterius/core"
	"github.com/Nyarum/noterius/database"
	"github.com/Nyarum/noterius/entitie"
	"github.com/Nyarum/noterius/manager"
	"github.com/Nyarum/noterius/packet"
	log "github.com/Sirupsen/logrus"
)

// ConnectHandler method for accept new connection from socket
func ConnectHandler(buffers *core.Buffers, app Application, c net.Conn) {
	defer core.ErrorNetworkHandler(c)

	var (
		buffer      *bytes.Buffer      = bytes.NewBuffer([]byte{})
		player      *entitie.Player    = entitie.NewPlayer(buffers)
		database    *database.Database = database.NewDatabase(&app.DatabaseInfo)
		manager     *manager.Manager   = manager.NewManager(database)
		packetAlloc *packet.Packet     = packet.NewPacket(player, manager)
	)

	// Once send first a packet
	packet, err := packetAlloc.Encode(940)
	if err != nil {
		log.WithError(err).Error("Error in packet encode")
	}

	buffers.GetWC() <- packet

	for getBytes := range buffers.GetRC() {
		buffer.Reset()
		buffer.Write(getBytes)

		log.WithField("bytes", fmt.Sprintf("% x", buffer.Bytes())).Info("Print message from client")

		// Ping <-> pong
		if buffer.Len() <= 2 {
			buffers.GetWC() <- []byte{0x00, 0x02}
			continue
		}

		opcodes, err := packetAlloc.Decode(buffer.Bytes())
		if err != nil {
			log.WithError(err).Error("Error in packet decode")
			return
		}

		if len(opcodes) == 0 {
			continue
		}

		for _, opcode := range opcodes {
			response, err := packetAlloc.Encode(opcode)
			if err != nil {
				log.WithError(err).Error("Error in packet encode")
				break
			}

			buffers.GetWC() <- response
		}
	}
}