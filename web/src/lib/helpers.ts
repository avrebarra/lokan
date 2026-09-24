// kan helpers — trimmed for lokan single-user (no S3/env)
export function inferInitialsFromEmail(email: string): string {
  const localPart = email.split('@')[0]
  if (!localPart) return ''
  const separators = /[._-]/
  const parts = localPart.split(separators)
  if (parts.length > 1) {
    return ((parts[0]?.charAt(0) ?? '') + (parts[parts.length - 1]?.charAt(0) ?? '')).toUpperCase()
  }
  return localPart.slice(0, 2).toUpperCase()
}

export function getInitialsFromName(name: string): string {
  return name
    .split(' ')
    .map((part) => part.charAt(0).toUpperCase())
    .join('')
}
