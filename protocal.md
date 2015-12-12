### Simple IM protocal with scalable header based on message type

#### scalable header
- msg_type u_int16
- src_id u_int16
- dst_id u_int16
- error u_int16
- len u_int16

#### payload
- msg_type == 0 # login
  - src_id = nil
  - dst_id = nil
  - len = len(payload)
  - payload = nickname (less than 16bytes)

- msg_type == 1 # login ack
 - src_id = nil
 - dst_id = increment_id # create by server
 - len = len(increment_id)
 - payload = increment_id

- msg_type == 2 # send msg to unique user
  - src_id = src_id
  - dst_id = dst_id
  - len = len(payload)
  - payload = msg_to_send

- msg_type == 3 # send broadcast msg
  - src_id = src_id
  - dst_id = nil
  - len = len(payload)
  - payload = msg_to_send

- msg_type == 4 # send msg ack
  - src_id = nil
  - dst_id = client_id
  - len = len(payload)
  - error = 1/2/3/4/5/6
  - payload = customized string depend on error
    - 1: server recv msg success
	- 2: user recv msg success
	- 3: permission denied
	- 4: failed with unknown error
	- 5/6: backup
