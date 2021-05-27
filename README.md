# desproxy
Go tooling for proxying and intercepting DESFire APDUs. It uses a generic PC/SC reader to interface with a DESFire card, and an ACR122u with an NXP PN532 controller to emulate the card on the other side. This allows you to present the emulated card to i.e. a fare gate, and log and potentially modify all of the traffic. This is in contrast to sniffing the direct communication of a card and reader, where frames can be lost or corrupted, and cannot be modified/dropped.

### Installation
`desproxy` is simply a Go CLI application, but it does not support passing arguments; to modify its behavior, you'll need to edit the source code. Only a few batteries are included. :)

```
git clone git@github.com:iangcarroll/desproxy.git
cd desproxy && go run .
```

### Usage
By default, `desproxy` uses the first reader it finds as the emulator, and the second reader as the target. It immediately begins an emulation session with the emulator, and waits for a reader to connect to the emulated card. An example log is like this:

```
2021/05/26 15:55:14 Connecting to reader ACS ACR122 0
2021/05/26 15:55:14 Connecting to reader OMNIKEY CardMan 5x21-CL 0
2021/05/26 15:55:14 Serialized EmuReq to 38 bytes
2021/05/26 15:55:15 Received 5a909090
2021/05/26 15:55:15 Sending the target: 5a909090
2021/05/26 15:55:15 Target responded 1 00 <nil>
2021/05/26 15:55:15 Sending 00 back.
```

In this log, the reader connected to the ACR122 and sent a native DESFire command `0x5a` (select application `0x909090`). The command was successful, so the card returned `00` (no error) and the response was passed to the emulated card, which then sent it to the reader.

While an unlimited amount of commands or attempts can be used with the emulator, care should be taken to understand the fragility of the PN532's emulation setup. It is easy for `desproxy` to become desynchronized from the true emulation state, and bugs have especially been noticable when a reader soft or hard resets the emulated card, which is not detected by `desproxy` and thus the state of the card can become corrupted.

For example, when Apple Pay connects to certain transit cards, it runs a `validateCardScript` which sends multiple commands to the card to verify it is the correct kind. It seems to issue a soft reset in between these commands, which prevents bugs when it calls commands like `0x60` that throw an aborted command error when they are not fully read out. However, you will likely need to handle these types of issues yourself -- adding logic to `desproxy.go` to call `coldResetCard` in certain circumstances should be pretty easy.

### Problems
I've had to abandon using the ACR122 for certain applications because it does not support emulating a 7-byte UID; only three of the four UID bytes are changable via its emulation API. As a result, if cards check or use key derivation based on the card UID, the proxy will prevent authentication from succeeding.

### Other work
These resources were very helpful in building this:
* https://salmg.net/2017/12/11/acr122upn532-nfc-card-emulation/
* https://github.com/AdamLaurie/RFIDIOt