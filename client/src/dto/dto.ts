export interface BasicResp<T> {
  Data: T
  Message: string
  Code: number
}


const serverItem = {
  "id": 3,
  "serverName": "ExpansionServer",
  "language": "node-http",
  "upStreamName": "up_expansion_server",
  "location": "/expansionserver/"
}

export type Item = typeof serverItem