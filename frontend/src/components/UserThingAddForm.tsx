import { JSX, createSignal } from 'solid-js';
import Input from '../ui/Input/Input';
import AttachFile from '../ui/AttachFile/AttachFile';
import ThingPhoto from '../ui/ThingPhoto/ThingPhoto';
import TextArea from '../ui/TextArea/TextArea';
import Select from '../ui/Select/Select';
import SubmitButton from '../ui/SubmitButton/SubmitButton';
import Form from '../ui/Form/Form';
import addThing from '../utils/addThing';
import { allSymbolsRegExp, allSymbolsRegExpStr } from '../utils/regExps';
import { ThingType } from '../utils/consts';
import BackButton from '../ui/BackButton/BackButton';
import FormTitle from '../ui/FormTitle/FormTitle';
import fileToBase64 from '../utils/fileToBase64';

const UserThingAddForm = (props: {
  defaultThingType: ThingType;
}): JSX.Element => {
  const [thingType, setThingType] = createSignal(props.defaultThingType);
  const [thingName, setThingName] = createSignal('');
  const [thingPhoto, setThingPhoto] = createSignal('');
  const [userMessage, setUserMessage] = createSignal('');

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();
    addThing({
      thing: {
        type: thingType(),
        name: thingName(),
        userMessage: userMessage(),
        photo: thingPhoto(),
      },
    });
  };

  return (
    <Form onsubmit={handleSubmit}>
      <BackButton />
      <FormTitle>Создание объявления</FormTitle>
      <Select
        id='thing_type_select'
        value={thingType()}
        oninput={(event) => {
          event.target.value === ThingType.found
            ? setThingType(ThingType.found)
            : setThingType(ThingType.lost);
        }}
        label='Что случилось?*'
        required
      >
        <option
          value={ThingType.lost}
          selected={
            props.defaultThingType === ThingType.lost ? true : undefined
          }
        >
          Я потерял вещь
        </option>
        <option
          value={ThingType.found}
          selected={
            props.defaultThingType === ThingType.found ? true : undefined
          }
        >
          Я нашёл вещь
        </option>
      </Select>
      <Input
        placeholder='Название вещи*'
        name='thing_name'
        value={thingName()}
        oninput={(event) => setThingName(event.target.value)}
        required
        pattern={allSymbolsRegExpStr}
      />
      <TextArea
        placeholder='Сообщение'
        name='user_message'
        value={userMessage()}
        oninput={(event) => {
          setUserMessage(event.target.value);
          if (userMessage() != '' && !allSymbolsRegExp.test(userMessage())) {
            event.target.setCustomValidity(
              'Введите данные в указанном формате.',
            );
          } else {
            event.target.setCustomValidity('');
          }
        }}
      />
      <AttachFile
        accept='image/jpeg,image/png'
        id='attach_thing_photo'
        label='Выберите фотографию'
        oninput={(event) =>
          event.target.files &&
          fileToBase64(event.target.files[0]).then((photoBase64) =>
            setThingPhoto(photoBase64),
          )
        }
      />
      <ThingPhoto
        src={thingPhoto()}
        deletePhoto={() => setThingPhoto('')}
      />
      <SubmitButton name='add_thing_submit'>Добавить вещь</SubmitButton>
    </Form>
  );
};

export default UserThingAddForm;
