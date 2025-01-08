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



export interface T_Grid {
  id: number
  servantId: number
  port: number
  nodeId: number
  status: number
  pid: number
  updateTime: string
  gridServant: GridServant
  gridNode: GridNode
}

export interface GridServant {
  servantId: number
  language: string
  servantGroupId: number
  serverName: string
  servantCreateTime: string
  execPath: string
  protocol: string
  preview: string
}

export interface GridNode {
  id: number
  ip: string
  main: string
  nodeStatus: number
  nodeCreateTime: string
}

export interface T_StatLogListItem {
  id: number
  serverName: string
  pid: number
  threads: number
  isRunning: boolean
  createTime: string
  stat: string
}

export interface T_RelaseItem {
  id: number
  serverName: string
  filePath: string
  updateTime: string
}

export interface T_ServerList {
  id: number
  tagName: string
  tagEnglishName: string
  servants: Servant[]
}

export interface Servant {
  id: number
  serverName: string
  createTime: string
  language: string
}

export interface T_StatusLogList {
  list: T_Status[]
  total: number
}

export interface T_Status {
  id: number
  gridId: number
  stat: string
  createTime: string
  isRunning: string
  gridInfo?: T_Grid
}
