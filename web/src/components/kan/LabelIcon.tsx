// kan LabelIcon — small dot for label/colour chip
export default function LabelIcon({ colourCode }: { colourCode: string | null }) {
  return (
    <svg fill={colourCode ?? '#3730a3'} className="h-2 w-2" viewBox="0 0 6 6" aria-hidden="true">
      <circle cx={3} cy={3} r={3} />
    </svg>
  )
}
