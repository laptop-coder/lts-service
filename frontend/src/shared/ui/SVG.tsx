import type { Component } from 'solid-js';
import type { SVG } from '../model/index';

export const SVG = (props: SVG) => {
  return (
    <svg
      class={props.class}
      xmlns='http://www.w3.org/2000/svg'
      viewBox='0 -960 960 960'
    >
      <path d={props.d} />
    </svg>
  );
};
