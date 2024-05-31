import _ from 'lodash'
export function isEmptyObj(target:any):boolean{
    if(_.isObject(target)){
        if(Object.keys(target).length == 0){
            return true
        }
        return false
    }
    return false
}