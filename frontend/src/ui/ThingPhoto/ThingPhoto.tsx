import { JSX, Show } from 'solid-js';

import styles from './ThingPhoto.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';

const ThingPhoto = (
  props: JSX.ImgHTMLAttributes<HTMLImageElement> & { deletePhoto?: Function },
): JSX.Element => (
  <Show when={props.src}>
    <div class={styles.thing_photo_wrapper}>
      {props.deletePhoto && (
        <button
          class={styles.delete_thing_photo_button}
          type='button'
          onclick={() => {
            if (
              confirm(
                'Подтвердите удаление фотографии. Это действие необратимо',
              ) &&
              props.deletePhoto !== undefined
            ) {
              props.deletePhoto();
            }
          }}
        >
          <img src={`${ASSETS_ROUTE}/delete.svg`} />
        </button>
      )}
      <img
        src={props.src}
        title={props.title}
        class={styles.thing_photo}
        onclick={(event) => event.target.requestFullscreen()}
      />
    </div>
  </Show>
);

export default ThingPhoto;
