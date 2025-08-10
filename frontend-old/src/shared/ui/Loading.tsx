import type { Component } from 'solid-js';
import { SVG } from './index';
import { d } from '../assets/index';

export const Loading: Component = () => {
  return (
    <SVG
      d={d.loading}
      class='loading'
    />
  );
};
