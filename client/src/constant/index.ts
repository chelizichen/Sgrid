export const constants = {
  TOKEN: 'access_token',
  LANGUAGES: [
    'node',
    'java',
    'java(jar)',
    'go',
    'exe',
    'custom command',
    'python(tar)',
    'python(exe)'
  ]
}

export function localGet(key: string) {
  return localStorage.getItem(key)
}

export function localSet(key: string, value: string) {
  return localStorage.setItem(key, value)
}

export function localDel(key: string) {
  return localStorage.removeItem(key)
}
