// AUTOMATICALLY GENERATED - DO NOT EDIT

export interface EchoRequest {
	"message": string;
}
export interface EchoResponse {
	"message": string;
}


export interface TimeResponse {
	"Time": number /* unix timestamp */;
}
export interface ExtendedField {
	"Test": string;
	"iamfloating": number /* float32 */;
	"Banana": boolean;
	"Flotarr": number /* float32 */[];
	"Bytes": Uint8Array;
}
export interface TestStruct {
	"Test": string;
	"iamfloating": number /* float32 */;
}

/// RPC Generated
export interface App {
	BinaryExample(): Promise<Uint8Array>
	Example(arg0: EchoRequest): Promise<EchoResponse>
	ExamplePlain(arg0: EchoRequest): Promise<string>
	FailForMe(arg0: EchoRequest): Promise<string>
	GetTime(): Promise<TimeResponse>
	NativeTypeExample(arg0: number /* float32 */[]): Promise<number /* float64 */>
	RandBytes(): Promise<Uint8Array>
	TypeHandling(): Promise<ExtendedField>
	TypeHandling2(): Promise<TestStruct>
}