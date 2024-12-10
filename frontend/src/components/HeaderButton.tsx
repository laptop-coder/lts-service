import type { Component } from 'solid-js';

import styles from './headerbutton.module.css';
import svg_styles from '../svg.module.css';


interface HeaderButtonProps {
  d: string,
  action: func
}


const HeaderButton: Component = (props: HeaderButtonProps) => {
  return (
    <button class={styles.header_button} onClick={props.action}>
      <svg class={svg_styles.svg} xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
        <path d={props.d} />
      </svg>
    </button>
  );
}


export default HeaderButton;
