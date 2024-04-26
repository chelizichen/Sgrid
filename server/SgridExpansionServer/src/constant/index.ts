export enum constant {
  SGRID_TARGET_PORT = "SGRID_TARGET_PORT",
  SGRID_SERVER_CONF = "SGRID_SERVER_CONF",
  SIMP_SERVER_STORAGE = "SIMP_SERVER_STORAGE",
  SIMP_TARGET_PORT = "SIMP_TARGET_PORT"
}
export enum tables {
  eff_event = "eff_event",
  eff_eve_comment = "eff_eve_comment",
  blog_user = "blog_user"
}

export enum eventStatus {
  created = 0,
  doing = 1,
  done = 2,
  overTime = 3,
  failed = -1
}

export enum dates {
  FMT = "YYYY-MM-DD HH:mm:ss"
}
export const NewError = function (code: number, msg: string) {
  return {
    code: code || -1,
    msg
  }
}

export function getConf() {
  return {
    nginxPath: process.env.nginxPath!,
    historyDir: process.env.historyDir!
  }
}
