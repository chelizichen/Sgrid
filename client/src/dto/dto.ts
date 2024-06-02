export interface BasicResp<T> {
  data: T
  message: string
  code: number
}


const serverItem = {
  "id": 3,
  "serverName": "ExpansionServer",
  "language": "node-http",
  "upStreamName": "up_expansion_server",
  "location": "/expansionserver/"
}

export type Item = typeof serverItem

export interface Tree {
  id: number;
  label: string;
  children?: Tree[];
}
export type MenuVo = {
  id: number;
  name: string;
  title: string;
  path: string;
  component: string;
  parentId: number;
};
