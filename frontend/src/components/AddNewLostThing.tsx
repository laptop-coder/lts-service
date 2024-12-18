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
  const [addNewFoundThing, setAddNewFoundThing] = createSignal(false);

  const [thingName, setThingName] = createSignal("");
  const [userContacts, setUserContacts] = createSignal("");
  const [thingLocation, setThingLocation] = createSignal("");
  const [userMessage, setUserMessage] = createSignal("");

  const [data, setData] = createSignal({});

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
	      onClick={() => {setChooseThingType(false); setAddNewFoundThing(true)}}
	      type="text"
	      value="Я нашёл вещь"
	    />
	  </>
	}
	{addNewLostThing() &&
	  <>
	    <p>Добавить потерянную вещь</p>
	    <form method="post">
	      <input
	        placeholder="Что Вы потеряли?"
		value={thingName()}
		onInput={e => setThingName(e.target.value)}
		required
	      />
	      <input
	        placeholder="Как с Вами можно связаться?"
		value={userContacts()}
		onInput={e => setUserContacts(e.target.value)}
		required
	      />
	      <textarea
	        placeholder="Здесь можно оставить сообщение"
		value={userMessage()}
		onInput={e => setUserMessage(e.target.value)}
		required
	      />
	      <button
		onClick={e => {
		  e.preventDefault();
		  setData(JSON.stringify({
		    "thingName": thingName(),
		    "userContacts": userContacts(),
		    "userMessage": userMessage(),
		  }))
		}}
	      >
	        Отправить
	      </button>
	    </form>
	  </>
	}
	{addNewFoundThing() &&
	  <>
	    <p>Добавить найденную вещь</p>
	    <form method="post">
	      <input
	        placeholder="Что Вы нашли?"
		value={thingName()}
		onInput={e => setThingName(e.target.value)}
		required
	      />
	      <input
	        placeholder="Где забрать вещь?"
		value={thingLocation()}
		onInput={e => setThingLocation(e.target.value)}
		required
	      />
	      <textarea
	        placeholder="Здесь можно оставить сообщение"
		value={userMessage()}
		onInput={e => setUserMessage(e.target.value)}
		required
	      />
	      <button
		onClick={e => {
		  e.preventDefault();
		  setData(JSON.stringify({
		    "thingName": thingName(),
		    "thingLocation": thingLocation(),
		    "userMessage": userMessage(),
		  }))
		}}
	      >
		Отправить
	      </button>
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

