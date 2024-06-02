const views = import.meta.glob('./views/**/*.vue',{
    eager:true,
})

export const loadView = (view:string) => {
    return ()=>{
        return new Promise((resolve)=>{
            resolve(views[view])
        })
    }
}