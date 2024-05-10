import path from "path"
import { errorHandler, initHistroyDir } from "./src/configuration"
import SgridController from "./src/routes/nginx"
import { NewSgridServer, NewSgridServerCtx } from "sgridnode/build/main"
import { f_env } from "sgridnode/build/lib/constant/index"
function boost() {
  const ctx = NewSgridServerCtx()
  const conf = ctx.get(f_env.ENV_SGRID_CONFIG) as SimpConf
  setEnv({ nginxPath: conf.config.ng_file, historyDir: conf.config.ng_dir })
  initHistroyDir()
  const servant = path.join("/", conf.server.name.toLowerCase())
  const sgridController = new SgridController(ctx)
  ctx.use(servant, sgridController.router!)
  ctx.use(errorHandler())
  NewSgridServer(ctx)
}

boost()

process.on("uncaughtException", (err) => {
  console.error(err)
})

process.on("unhandledRejection", (reason, p) => {
  console.error(reason, p)
})

function setEnv(obj: Record<string, string>) {
  const keys = Object.keys(obj)
  for (const key of keys) {
    process.env[key] = obj[key]
  }
}
