import type { Component } from 'solid-js';
import { createSignal } from 'solid-js';

import Button from './Button';
import styles from './addnewlostthing.module.css';
import button_styles from './button.module.css';


interface AddNewLostThingProps {
  onClick: func;
}

const AddNewLostThing: Component = (props: AddNewLostThingProps) => {
  const [chooseThingType, setChooseThingType] = createSignal(true);
  const [addNewLostThing, setAddNewLostThing] = createSignal(false);
  return (
    <div class={styles.wrapper}>
      <div class={styles.box}>
	{chooseThingType() &&
	  <>
	    <Button
	      class={button_styles.wide_button}
	      onClick={() => {setChooseThingType(false); setAddNewLostThing(true)}}
	      type="text"
	      value="Я потерял вещь"
	    />
	    <Button
	      class={button_styles.wide_button}
	      onClick={() => console.log("Found thing")}
	      type="text"
	      value="Я нашёл вещь"
	    />
	  </>
	}
	{addNewLostThing() &&
	  <>
	    <p>Добавить потерянную вещь</p>
	    <form method="post">
	      <input placeholder="Что Вы потеряли?" required />
	      <input placeholder="Как с Вами можно связаться?" required />
	      <textarea placeholder="Здесь можно оставить сообщение" required />
	      <button type="submit">Отправить</button>
	    </form>
	  </>
	}
      </div>
      <div class={styles.background} onClick={props.onClick}>
      </div>
    </div>
  );
}


export default AddNewLostThing;

