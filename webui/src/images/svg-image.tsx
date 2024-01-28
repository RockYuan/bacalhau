// Write a component that renders an SVG image. The component should accept a
// `src` prop that specifies the URL of the image to render. The component
// should also accept an `alt` prop that specifies the alternative text for the
// image. The component should render the image and the alternative text inside
// an SVG element. The SVG element should have a `role` attribute with the value
// `"img"`. The SVG element should also have an `aria-label` attribute with the
// value of the `alt` prop.

// The component should render the following SVG element:
// <svg role="img" aria-label="{alt}">
//   <image href="{src}" />
// </svg>

// Reusable SVG rendering component below:
import React from "react"

interface SVGImageProps {
  src: string
  alt: string
}

export const SVGImage: React.FC<SVGImageProps> = ({ src, alt }) => (
  // use alt prop to set aria-label
  <svg role="img" aria-label={alt}>
    <image href={src} />
  </svg>
)
