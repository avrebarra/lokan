// vendored from kan — https://github.com/kan/kan (AGPLv3)
// source: kan/apps/web/src/views/board/dnd/types.ts — see NOTICE.md
export interface ListDragData {
  type: 'LIST'
}

export interface ListBodyDragData {
  type: 'LIST_BODY'
  listPublicId: string
}

export interface CardDragData {
  type: 'CARD'
  listPublicId: string
}

export type DragData = ListDragData | ListBodyDragData | CardDragData
