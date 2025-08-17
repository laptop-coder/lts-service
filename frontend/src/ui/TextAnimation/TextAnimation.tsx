import { JSX, ParentProps } from 'solid-js';

import { Motion } from 'solid-motionone';

const TextAnimation = (props: ParentProps): JSX.Element => (
  <Motion.span
    animate={{ opacity: [0, 1] }}
    transition={{ duration: 0.5, easing: 'ease-in-out' }}
  >
    {props.children}
  </Motion.span>
);

export default TextAnimation;
