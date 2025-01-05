import type { Component } from "solid-js";

interface SVGProps {
  d: string;
}

const SVG = (props: SVGProps) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 -960 960 960"
    >
      <path d={props.d} />
    </svg>
  );
};

export default SVG;
