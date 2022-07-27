package modbus

import (
    "io"
    "time"
)

type TCPRtuClientHandler struct {
    rtuPackager
    tcpTransporter
}


func (mb *TCPRtuClientHandler) Send(aduRequest []byte) (aduResponse []byte, err error) {
    mb.mu.Lock()
    defer mb.mu.Unlock()

    // Establish a new connection if not connected
    if err = mb.connect(); err != nil {
        return
    }
    // Set timer to close when idle
    mb.lastActivity = time.Now()
    mb.startCloseTimer()
    // Set write and read timeout
    var timeout time.Time
    if mb.Timeout > 0 {
        timeout = mb.lastActivity.Add(mb.Timeout)
    }
    if err = mb.conn.SetDeadline(timeout); err != nil {
        return
    }
    // Send data
    mb.logf("modbus: sending % x", aduRequest)
    if _, err = mb.conn.Write(aduRequest); err != nil {
        return
    }
    function := aduRequest[1]
    functionFail := aduRequest[1] & 0x80
    bytesToRead := calculateResponseLength(aduRequest)

    time.Sleep(time.Millisecond * 100)

    var n int
    var n1 int
    var data [rtuMaxSize]byte
    //We first read the minimum length and then read either the full package
    //or the error package, depending on the error status (byte 2 of the response)
    n, err = io.ReadAtLeast(mb.conn, data[:], rtuMinSize)
    if err != nil {
        return
    }
    //if the function is correct
    if data[1] == function {
        //we read the rest of the bytes
        if n < bytesToRead {
            if bytesToRead > rtuMinSize && bytesToRead <= rtuMaxSize {
                if bytesToRead > n {
                    n1, err = io.ReadFull(mb.conn, data[n:bytesToRead])
                    n += n1
                }
            }
        }
    } else if data[1] == functionFail {
        //for error we need to read 5 bytes
        if n < rtuExceptionSize {
            n1, err = io.ReadFull(mb.conn, data[n:rtuExceptionSize])
        }
        n += n1
    }

    if err != nil {
        return
    }
    aduResponse = data[:n]
    mb.logf("modbus: received % x\n", aduResponse)
    return
}


func NewTCPClient2Handler(address string) *TCPRtuClientHandler {
    h := &TCPRtuClientHandler{}
    h.Address = address
    h.Timeout = tcpTimeout
    h.IdleTimeout = tcpIdleTimeout
    return h
}
