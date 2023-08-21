## Connection pool
To allow the support of concurrent requests and responses from the merchant to the payment network, we are using a
connection pool.

### Components of the connection pool system
#### The pool
* Starts the connections when there is a `Start` call.
* Accepts new requests to the system and passes them to free connections if any. When the connections are busy managing
a request, they don't listen to new pool requests until they have finished with the current request.
* Shuts down the connections when there is a `Stop` call.

#### The connection
This component is in charge of:
* Maintaining a connected socket with the payment network.
* Accept new request jobs whenever the connection is connected to the payment network and not attending other job.
* Identify and notify of the response to a request.
* Notify of payment network request.
* Receive and send packets over the network.

The connection is managing a socket that is used concurrently for receiving, sending and connecting. Trying to keep
alive the socket session against the payment network, sending packets to the payment network and receiving packets from
the payment network.<br>
To achieve that, the selected method has been channels and goroutines instead of interlocking and waits.
Each goroutine is in charge of only one thing. Each goroutine except the `run`one, are state machines: each one adopts
a behavior (selected by the method `selectxxxxBehavior`where `xxxx` is the name of the goroutine) depending on the
current goroutine state. In that behavior the goroutine only attends to the channels that can be applied to that state,
for example: the `attend` goroutine (the one that is in charge of accepting new requests from the merchant) only attends
to the new request signal if it is connected and idle (not attending any other request).<br>
The channels that synchronize the goroutines are:
* `shutdownSignal`: This channel is signaled by the pool on the `stop` method. All the goroutines and connections should
  stop.
*
The goroutines that make the connection work are:
* [`connect`](./connect.go)
  * Maintains a connected socket to the payment network. If the connection breaks, it tries to reconnect the socket on 
  its own every `reconnectDelay`
    * It has two behaviors depending on the connected state:
      * `disconnectedConnectorBehavior`
        * Applied when the socket is closed and the component is disconnected from the payment network. This behavior 
        only attends to the following channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `recontectTicker`: This ticker notifies to the goroutine that it should try a new connection attempt.
            * If the connection attempt has been successful, it signals `recieve`, `send` and `attend` goroutines that 
            the socket is connected and marks the current one with the connected state.
          * `disconnectCmdSignal`: Just in case, but not required. This channel signals whenever any other goroutine 
          wants to close the connection.
      * `connectedConnectorBehavior`
        * Applied when the socket has been connected to the payment network. This behavior attends to the following 
        channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `disconnectCmdSignal`: Closes the current socket and starts the reconnection ticker. This channel is 
          signaled by other goroutines. It makes the `connect` enter the disconnected state and signal `receive`, `send`
          and `attend`goroutines that the socket has been disconnected.
* [`send`](./send.go)
  * As there can not be concurrent sends on a socket as they would mix bytes in the socket buffer, all the sends should 
  be sequential. This is achieved by setting up this goroutine that `receive` and `attend` use to send packets over the 
  socket. It is also a state machine.
    * It has two behaviors depending on the connected state:
      * `disconnectedSenderBehavior`
        * Applied when the socket is closed and the component is disconnected from the payment network. This behavior 
        attends to the following channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `senderConnectedSignal`: This channel is signaled by the `connect`goroutine whenever the socket is up and 
          connected. It makes the `send` goroutine enter the connected state.
      * `connectedSenderBehavior`
        * Applied when the socket is connected to the payment network. Attends to the following channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `senderDisconnectedSignal`: To enter the disconnected state. This channel is signaled by the `connect` 
          goroutine.
          * `sendCmdSignal`: To attend the send command requests. Each send command is composed by a buffer to be sent 
          and a flag indicating that this goroutine should notify of the error if any using the `errorSignal`. In case 
          of error it will notify the error if required, notify `connect` that it wants to close the socket and enters 
          the disconnected state.
* [`receive`](./receive.go)
  * This goroutine is in charge of receiving packets from the network and notify of the reception. It is also in charge 
  of notifying the `send` goroutine to send the zero length probe if it doesn't receive anything for a `keepAliveDelay`.
  This is done to avoid the other side to close the connection due to inactivity. As the previous mentioned goroutines 
  it is also a state machine.
    * It has two behaviors depending on the goroutine status:
      * `disconnectedReceiverBehavior`
        * Applied when the socket is closed and the component is disconnected from the payment network. This behavior 
        attends to the following channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `receiverConnectedSignal`: To begin receiving over the socket. It makes the goroutine enter the connected 
          state.
      * `connectedReceiverBehavior`
        * Applied when the socket is connected to the payment network. Attends to the following channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `receiverDisconnectedSignal`: Signaled by the `connect`goroutine whenever the socket is disconnected. Makes 
          this goroutine enter the disconnected state.
          * If any of the previous signals would block (they are not signaled). It tries to receive a packet from the 
          network.
            * No packet is received before `keepAliveDelay`: It notifies `send` goroutine, via `sendCmdSignal`to send a
            zero byte probe to keep the socket alive and connected.
            * A packet is received: notifies the `attend` goroutine that a packet has been received via `receivedSignal`
            .
            * An error occurred in the reception: notifies the error via `errorSignal` and notifies `connect` that it 
            should disconnect the socket via `disconnectCmdSignal`. The goroutine enters in disconnected state.
* [`attend`](./attend.go)
  * This goroutine is in charge of managing all the requests/responses over the connection component. It is also a state
  machine.
    * It has three behaviors depending on the connected state and if the goroutine has a request pending or not.
      * `disconnectedAttendantBehavior`
        * Applied when the socket is closed and the component is disconnected from the payment network. This behavior
        attends to the following channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `attendantConnectedSignal`: Signaled by the `connect` go routine, it enters the connected and idle state.
      * `connectedIdleAttendantBehavior`
        * Applied whenever the goroutine is in connected state, but hasn't got any request pending. It attends the
        following channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `attendantDisconnectedSignal`: Signaled by the `connect` goroutine. It makes to enter disconnected state.
          * `requestSignal`: Signaled by the external code via pool's send method. It stores the request as a pending
          request, notifies the `send` goroutine that it has to send the request via the `sendCmdSignal` and enters the
          busy state.
          * `receivedSignal`: Signaled by the `receive` goroutine where there is a network packet. In this state the
          packets that we should receive are network management packets from the payment network. Uses the connection
          `receivedFactory` to decode the received package. Connection's `receivedFactory` is a function which
          implementation is payment network dependant so this is a dependency that the developer should pass into the 
          component. Once the packet has been decoded it tries to get a response packet using `PacketToSend` method.
          This method is payment network dependant so this is a dependency the developer should implement. If there is a
          packet to send, it notifies the `sender` goroutine via `sendCmdSignal`.
      * `connectedBusyAttendantBehavior`
        * Applied if the goroutine is in connected state and has a pending request (is waiting for a response). It 
        attends to the following channels:
          * `shutdownSignal`: To properly shut down the goroutine.
          * `attendantDisconnectedSignal`: To enter disconnected state. It finishes the pending request with
          `ErrAbandoned`.
          * `errorSignal`: Signaled by `receive`and `send` goroutines. It finishes the pending request with the error
          signaled.
          * `receivedSignal`: Signaled by `receive`. It uses connection's `receivedFactory` as explained before. 
            * Packet is a response for the pending request: Finish the pending request with the received package.
            * Packet is not a response for the pending request: Behaves as explained in `connectedIdleAttendantBehavior`
            .
      * `pendingJob.context.Done`
        * Signaled if the context of the pending job has been canceled. It finishes the pending job with an 
        `ErrRequestCanceled`.
      * `time.After(c.responseTimeout)`
        * Signaled if the timeout waiting a pending job response has passed. It finishes the pending job with an 
        `ErrRequestTimeout`. 
* [`run`](./connection.go)
  * This goroutine is not a state machine and is in charge of starting the other goroutines and wait for them to finish
  after the `shutdownSignal` is signaled.

#### Messages
To communicate with the payment network, the pool uses two interfaces and a function to work with requests and 
responses. The implementation of the interfaces and the factory function is payment network dependent.

#####`Request` interface
The `Request` interface defines the methods that the connection need to perform the request.
* `Packet`: This method encodes the payment network message to a byte buffer. Is used by the 
`attendant` go routine to get the request packet that is going to send via `send` goroutine.
For an example on how to implement this interface take a look at 
[`Request`](../../processing/scheme/mastercard/request.go)

#####`Received` interface
The `Received` interface defines the methods that connection need to manage received packets from the payment network.
* `IsRequestResponse`: This method determines if the `Received`packet is the response of the passed `Request`. It is 
used by the `attendant` goroutine when it is busy attending one request, to check if the received packet is the response
of the pending request.
* `PacketToSend`: This method returns the response, if any, to a network request. Is used by the `attendant` goroutine
when it receives something from the `receive` goroutine, and it is not a request's response.
* `Error`: This method returns the error stored in the `Received` interface.
For an example on how to implement this interface take a look at 
[`Response`](../../processing/scheme/mastercard/response.go)

####`ReceivedFactoryFunc` function
As the component does not know how to decode the received packets, it transforms the received packets into a `Received`
compliant object using this function. This function is payment network dependant. It is used by the `attend` goroutine 
to generate the `Received` object it requires.
For an example on how to implement this interface take a look at 
[`Response`](../../processing/scheme/mastercard/response.go). The function `NewMastercardResponse` would be the 
`ReceivedFactoryFunc`

### Connection lifecycle explanation

#### Connecting phase:
In this phase, only the `connect` goroutine is active trying to connect to 
the payment network using the reconnect ticker. When a successful connection has been made, all the goroutines are 
signaled that they should commute to connected. The `connect` goroutine becomes inactive waiting for disconnection 
commands.

#### Connected phase:
In this phase the `receive` goroutine is trying to receive packets from the network, the `send` goroutine is waiting 
send commands, the `connect` goroutine is waiting for disconnect commands and the `attend` goroutine is waiting for 
sending requests from the pool.

* A new sending request arrives from the pool
  * The `attend` goroutine enters he busy state, so it won't accept new sending requests, it sends the send command to 
  `send` goroutine with the payload it has to send and stores the sending request as pending.
    * If `send` fails, the pending sending request is cancelled with the error `send` had and `connect` is notified 
    that the connection should close the socket.
  * When `receive` goroutine receives a packet from the network, it notifies to the `attend` goroutine that a packet has
  been received. The `attend` goroutine will decode the packet received, will check if the received packet is a response
  for the pending request and if not tries to get a response for the received packet as it should be a network request. 
  If the received packet is a response it will notify the initial caller (pool.send method) that it has a response. If
  the received packet is not a response and it has a response, it will notify the `sender` goroutine that has something
  to send, and will continue waiting for another packets.
      * If `receive` has an error, it notifies the `attend` goroutine that if it has a pending request will notify to the 
    initial caller of the pool.Send. It will also notify `connect` goroutine that the socket should be closed.

  





