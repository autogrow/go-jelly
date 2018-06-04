// Package ig provides a client to the IntelliGrow API which allows users to interface programmatically with
// their IntelliClimate and IntelliDose devices.
//
// To connect to IntelliGrow a username and password must be provided on instantiation of the client:
//
//     client, err := ig.NewClient("me", "secret")
//     if err != nil {
//       panic(err)
//     }
//
// From there the client can be used to query the devices and growrooms attached to their account.  Readings
// can be requested and various actions can be taken to change settings of the device and even force dosing.
package ig
