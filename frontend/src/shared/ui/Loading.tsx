import type { Component } from 'solid-js';
import { d } from '../assets/d';
import { SVG } from '../ui/SVG';

export const Loading: Component = () => {
  return (
    <SVG
      d={d.loading}
      class='loading'
    />
  );
};
