
 [Install](#install-skdc)  | [Go Docs](https://pkg.go.dev/github.com/superkabuki/skdc)  | [Examples](https://pkg.go.dev/github.com/superkabuki/skdc) | 

# `S`uper `K`arate `D`eath `C`ar is a SCTE-35 Parser lib written in Go.
# Encoder/Decoder for SCTE-35
<br>


- [x] Parses SCTE-35 Cues from MPEGTS or Bytes or Base64 or Hex or Int or Octal or even Base 36.
- [x] Parses SCTE-35 Cues spread over multiple MPEGTS packets  
- [x] Supports multi-packet PAT and PMT tables  
- [x] Supports multiple MPEGTS Programs and multiple SCTE-35 streams 
- [x] Encodes Time Signals and Splice Inserts with Descriptors and Upids. 
 
___

##  Super Karate Death Car is fast.
![image](https://github.com/user-attachments/assets/edc29afe-d198-44e2-a230-567acbcddbdf)




### Want to parse an MPEGTS video and print the SCTE-35?  🛰️

### Do it in ten lines.
```go
package main                        

import (                              
        "os"                            
        "github.com/superkabuki/skdc"       
)                                    

func main(){                         
        arg := os.Args[1]             
        stream := skdc.NewStream()    
        stream.Decode(arg)           
}                                    
```
---

# Documentation


* [Install](#install-skdc)  


* [Examples](https://pkg.go.dev/github.com/superkabuki/skdc)
	* skdc.Stream 
		* [Parse SCTE-35 from MPEGTS](#quick-demo)
   		* [Custom Cue Handling for MPEGTS Streams](#custom-cue-handling-for-mpegts-streams)
   		* [Multicast](#custom-cue-handling-for-mpegts-streams-over-multicast) _(New!)_

   	* skdc.Cue
		* [Parse Base64 encoded SCTE-35](#parse-base64-encoded-scte-35) 
		* [Use Dot Notation to access SCTE-35 Cue values](#use-dot-notation-to-access-scte-35-cue-values)
		* [Shadow a Cue Struct Method ( override ) ](#shadow-a-cue-struct-method)
		* [Shadow a Cue Method and call the Shadowed Method ( like super in python )](#call-a-shadowed-method)
  		* [Load a SCTE-35 Cue from JSON and Encode it](#load-json-and-encode)



### `Install skdc` 

```go
go install github.com/superkabuki/skdc@latest

```



### `Quick Demo` 

* skdcdemo.go

```go
package main

import (
        "os"
        "fmt"
        "github.com/superkabuki/skdc"
)

func main(){

        arg := os.Args[1]

        stream := skdc.NewStream()
        cues := stream.Decode(arg)
        for _,cue := range cues {
        fmt.Printf("Command is a %v\n", cue.Command.Name)
        }

}



```

#### `build skdcdemo`
```go
go build skdcdemo.go
```
#### `parse mpegts video for scte35` 
```go
./skdcdemo a_video_with_scte35.ts
```
			
### SCTE-35 Cues
```go
type Cue struct {
	InfoSection *InfoSection  
	Command     *Command          
	Dll         uint16       `json:"DescriptorLoopLength"`
	Descriptors []Descriptor `json:",omitempty"`
	Crc32       string
	PacketData  *packetData `json:",omitempty"`
}
```
* Output is JSON and looks like this.
```awk
{
    "InfoSection": {
        "Name": "Splice Info Section",
        "TableID": "0xfc",
        "SectionSyntaxIndicator": false,
        "Private": false,
        "SapType": 3,
        "SapDetails": "No Sap Type",
        "SectionLength": 44,
        "ProtocolVersion": 0,
        "EncryptedPacket": false,
        "EncryptionAlgorithm": 0,
        "PtsAdjustment": 2.3,
        "CwIndex": "0x0",
        "Tier": "0xfff",
        "CommandLength": 10,
        "CommandType": 5
    },

    "Command": {
        "Name": "Splice Insert",
        "CommandType": 5,
        "SpliceEventID": 1,
        "SpliceEventCancelIndicator": false,
        "OutOfNetworkIndicator": false,
        "ProgramSpliceFlag": true,
        "DurationFlag": false,
        "BreakDuration": 0,
        "BreakAutoReturn": false,
        "SpliceImmediateFlag": true,
        "EventIDComplianceFlag": true,
        "UniqueProgramID": 39321,
        "AvailNum": 1,
        "AvailExpected": 1
    },
    "DescriptorLoopLength": 17,

    "Descriptors": [
        {
            "Tag": 2,
            "Length": 15,
            "Name": "Segmentation Descriptor",
            "Identifier": "CUEI",
            "SegmentationEventID": "0x0",
            "SegmentationEventCancelIndicator": false,
            "SegmentationEventIDComplianceIndicator": true,
            "ProgramSegmentationFlag": true,
            "SegmentationDurationFlag": false,
            "DeliveryNotRestrictedFlag": false,
            "WebDeliveryAllowedFlag": false,
            "NoRegionalBlackoutFlag": false,
            "ArchiveAllowedFlag": false,
            "DeviceRestrictions": "Restrict Group 0",
            "SegmentationDuration": 0,
            "SegmentationMessage": "Provider Placement Opportunity End",
            "SegmentationUpidType": 1,
            "SegmentationUpidLength": 0,
            "SegmentationUpid": null,
            "SegmentationTypeID": 53,
            "SegmentNum": 0,
            "SegmentsExpected": 0,
            "SubSegmentNum": 0,
            "SubSegmentsExpected": 0
        }
    ],
    "Crc32": " 0x2d974195",
    "PacketData": {
        "Pid": 258,
        "Program": 1,
        "Pts": 72951.196988
    }
}
```

### `Parse base64 encoded SCTE-35`
```go
package main

import (
	"fmt"
	"github.com/superkabuki/skdc"
)

func main(){

	cue := skdc.NewCue()
	data := "/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM="
        cue.Decode(data) 
        fmt.Println("Cue as Json")
        cue.Show()
}
```

### `Shadow a Cue struct method`
```go
package main

import (
	"fmt"
	"github.com/superkabuki/skdc"
)

type Cue2 struct {
    skdc.Cue               		// Embed skdc.Cue
}
func (cue2 *Cue2) Show() {        	// Override Show
	fmt.Printf("%+v",cue2.Command)
}

func main(){
	var cue2 Cue2
	data := "/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM="
        cue2.Decode(data) 
        cue2.Show()
}

```

### `Call a shadowed method`
```go
package main

import (
	"fmt"
	"github.com/superkabuki/skdc"
)

type Cue2 struct {
    skdc.Cue               		// Embed skdc.Cue
}
func (cue2 *Cue2) Show() {        	// Override Show
	fmt.Println("Cue2.Show()")
	fmt.Printf("%+v",cue2.Command) 
	fmt.Println("\n\nskdc.Cue.Show() from cue2.Show()")
	cue2.Cue.Show()			// Call the Show method from embedded skdc.Cue
}

func main(){
	var cue2 Cue2
	data := "/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM="
        cue2.Decode(data) 
        cue2.Show()
	
}


```
### `Use Dot notation to access SCTE-35 Cue values`
* SuperKarate DeathCar doesn't use interfaces for SCTE-35 data
* SCTE-35 Data can be accessed directly via dot notation 

```go

/**
Show  the packet PTS time and Splice Command Name of SCTE-35 Cues
in a MPEGTS stream.
**/
package main

import (
	"os"
	"fmt"
	"github.com/superkabuki/skdc"
)

func main() {
	arg := os.Args[1]
	stream := skdc.NewStream()
	cues :=	stream.Decode(arg)
	for _,c := range cues {
		fmt.Printf("PTS: %v, Splice Command: %v\n",c.PacketData.Pts, c.Command.Name )
	}
}

```

### `Load JSON and Encode`
* skdc can accept SCTE-35 data as JSON and encode it to Base64, Bytes, or Hex string.
* The function __skdc.Json2Cue()__ accepts SCTE-35 JSON as input and returns a *skdc.Cue

```go
package main

import (
	"fmt"
	"github.com/superkabuki/skdc"
)

func main() {

	js := `{
    "InfoSection": {
        "Name": "Splice Info Section",
        "TableID": "0xfc",
        "SectionSyntaxIndicator": false,
        "Private": false,
        "Reserved": "0x3",
        "SectionLength": 42,
        "ProtocolVersion": 0,
        "EncryptedPacket": false,
        "EncryptionAlgorithm": 0,
        "PtsAdjustment": 0,
        "CwIndex": "0xff",
        "Tier": "0xfff",
        "CommandLength": 15,
        "CommandType": 5
    },
    "Command": {
        "Name": "Splice Insert",
        "CommandType": 5,
        "SpliceEventID": 5690,
        "OutOfNetworkIndicator": true,
        "ProgramSpliceFlag": true,
        "TimeSpecifiedFlag": true,
        "PTS": 23683.480033
    },
    "DescriptorLoopLength": 10,
    "Descriptors": [
        {
            "Length": 8,
            "Identifier": "CUEI",
            "Name": "Avail Descriptor"
        }
    ],
    "Crc32": "0xd7165c79"
}
`
cue :=  skdc.Json2Cue(js)    // 
cue.AdjustPts(28.0)   	 // Apply pts adjustment
fmt.Println("\nBytes:\n\t", cue.Encode())	// Bytes
fmt.Println("\nBase64:\n\t",cue.Encode2B64())  	// Base64
fmt.Println("\nHex:\n\t",cue.Encode2Hex()) 	// Hex

}


```
* Output
```smalltalk
Bytes:
	[252 48 42 0 0 0 38 115 192 255 255 240 15 5 0 0 22 58 127 207 254 127 12 79 115
	0 0 0 0 0 10 0 8 67 85 69 73 0 0 0 0 236 139 53 78]

Base64:
	 /DAqAAAAJnPA///wDwUAABY6f8/+fwxPcwAAAAAACgAIQ1VFSQAAAADsizVO

Hex:
	 0xfc302a0000002673c0fffff00f050000163a7fcffe7f0c4f7300000000000a00084355454900000000ec8b354e

```
## skdc.Stream
### `Custom Cue Handling for MPEGTS Streams`
##### Four Steps
1) Create Stream Instance
2) Read Bytes from the video stream (__in multiples of 188__)
3) Call Stream.DecodeBytes(Bytes) 
4) Process [] *Cue returned by Stream.DecodeBytes

```go
package main

import (
        "fmt"
        "github.com/superkabuki/skdc"
        "os"
)

func main() {

  arg := os.Args[1]
  stream := skdc.NewStream()  //   (1)
  bufSize := 32768 * 188   // Always read in multiples of 188
  file, err := os.Open(arg)
  if err != nil {
    fmt.Printf("Unable to read %v\n", arg)
  }
  buffer := make([]byte, bufSize)
  for {
  	_, err := file.Read(buffer)   //  (2)
	if err != nil {
		break
	}
	cues := stream.DecodeBytes(buffer)  // (3)

	for _, c := range cues {  //   (4)
		fmt.Printf(" %v, %v\n", c.PacketData.Pts, c.Encode2B64())
	}
  }
}

```
* Output
```php
60638.745877, /DAWAAAAAAAAAP/wBQb/RUqw1AAAd6OnQA==
 60638.745877, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60640.714511, /DAWAAAAAAAAAP/wBQb/RU1wqAAAoqaOaA==
 60640.714511, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60642.015811, /DAWAAAAAAAAAP/wBQb/RU9F4AAA9Te5ag==
 60642.015811, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60642.749877, /DAWAAAAAAAAAP/wBQb/RVAwfAAAWOLrFQ==
 60642.749877, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60644.718511, /DAWAAAAAAAAAP/wBQb/RVLwUAAAj7/Pgw==
 60644.718511, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60646.720511, /DAWAAAAAAAAAP/wBQb/RVWwJAAA8pm/jg==
 60646.720511, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60648.121911, /DAWAAAAAAAAAP/wBQb/RVec0gAAt0QzqA==
 60648.121911, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60634.208011, /DAWAAAAAAAAAP/wBQb/RUR1fAAAik8gfQ==
 60634.208011, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60634.675144, /DAWAAAAAAAAAP/wBQb/RUUxLAAABQVyEA==
 60634.675144, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
 60636.710511, /DAWAAAAAAAAAP/wBQb/RUfxAAAA0lhWhg==
 60636.710511, /DAgAAAAAAAAAP/wDwUAAAABf//+AFJlwAABAAAAAMOOklg=
```

### `Custom Cue Handling for MPEGTS Streams Over Multicast`
##### Need a multicast sender? Try [gums](https://github.com/superkabuki/gums)

<div> for multicast we use the same four steps,<br> the only difference is we read the bytes from the network instead of a local file. </div>

1) Create Stream Instance
2) Read Bytes from the video stream (__in multiples of 188__)
3) Call Stream.DecodeBytes(Bytes) 
4) Process [] *Cue returned by Stream.DecodeBytes
```go
package main

import (
	"fmt"
	"github.com/superkabuki/skdc"
	"os"
        "net"
)


func main() {
  arg := os.Args[1]
  stream := skdc.NewStream()  //  (1)
  stream.Quiet = true
  dgram:=1316  // <-- multicast dgram size is 1316 (188*7) for mpegts
  bufSize := 100 * dgram
  addr, _ := net.ResolveUDPAddr("udp", arg)
  l, _ := net.ListenMulticastUDP("udp", nil, addr)  // Multicast Connection 
  l.SetReadBuffer(bufSize)
  for {
	buffer := make([]byte, bufSize)  // (2)
	l.ReadFromUDP(buffer)
	cues := stream.DecodeBytes(buffer)   //  (3)
	for _, c := range cues {      //  (4)
		fmt.Printf(" %v, %v\n", c.PacketData.Pts, c.Encode2B64())
	}
  }
}
```
