import type { Component } from 'solid-js';

import styles from './button.module.css';
import svg_styles from '../svg.module.css';


interface ButtonProps {
  class?: string;
  onClick: func;
  type: string;
  value: string;
}


const Button: Component = (props: ButtonProps) => {
  return (
    <button class={styles.button + (props.class ? " " + props.class : "")} onClick={props.onClick}>
      {props.type === "svg" && 
        <svg class={svg_styles.svg} xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
          <path d={props.value} />
        </svg>
      }
      {props.type === "text" &&
	props.value
      }
    </button>
  );
}

export default Button;

