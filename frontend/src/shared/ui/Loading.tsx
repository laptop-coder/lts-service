import type { Component } from 'solid-js';
import { SVG } from '../ui/SVG';
import { d } from '../assets/d';

export const Loading: Component = () => {
  return (
    <SVG
      d={d.loading}
      class='loading'
    />
  );
};
