import { JSX, createSignal, createResource, createEffect } from 'solid-js';
import Input from '../ui/Input/Input';
import AttachFile from '../ui/AttachFile/AttachFile';
import ThingPhoto from '../ui/ThingPhoto/ThingPhoto';
import TextArea from '../ui/TextArea/TextArea';
import Select from '../ui/Select/Select';
import SubmitButton from '../ui/SubmitButton/SubmitButton';
import Form from '../ui/Form/Form';
import editThing from '../utils/editThing';
import { allSymbolsRegExp, allSymbolsRegExpStr } from '../utils/regExps';
import { ThingType, STORAGE_ROUTE } from '../utils/consts';
import BackButton from '../ui/BackButton/BackButton';
import FormTitle from '../ui/FormTitle/FormTitle';
import fileToBase64 from '../utils/fileToBase64';
import checkPhotoAvailability from '../utils/checkPhotoAvailability';
import type { ResourceReturn } from 'solid-js';
import type { Thing } from '../types/thing';
import getThingData from '../utils/getThingData';
import deleteThingPhoto from '../utils/deleteThingPhoto';

const UserThingEditForm = (props: { thingId: string }): JSX.Element => {
  const [thingPhotoIsAvailable, setThingPhotoIsAvailable] = createSignal(false);
  const pathToPhoto = `${STORAGE_ROUTE}/${props.thingId}.jpeg`;
  checkPhotoAvailability({
    pathToPhoto: pathToPhoto,
    success: () => setThingPhotoIsAvailable(true),
  });
  const [oldThingData]: ResourceReturn<Thing> = createResource(
    { thingId: props.thingId },
    getThingData,
  );

  const [thingType, setThingType] = createSignal<ThingType>(ThingType.lost); // use "lost" as default value
  const [thingName, setThingName] = createSignal('');
  const [thingPhoto, setThingPhoto] = createSignal('');
  const [userMessage, setUserMessage] = createSignal('');

  let effectFirstRun = true;
  createEffect(() => {
    if (effectFirstRun && oldThingData.state === 'ready') {
      setThingType(oldThingData().Type);
      setThingName(oldThingData().Name);
      setUserMessage(oldThingData().UserMessage);
      setThingPhoto(oldThingData().Photo);
      effectFirstRun = false;
    }
  });

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();
    editThing({
      thing: {
        id: props.thingId,
        newType: thingType(),
        newName: thingName(),
        newUserMessage: userMessage(),
        newPhoto: thingPhoto(),
      },
    });
  };

  return (
    <Form onsubmit={handleSubmit}>
      <BackButton />
      <FormTitle>Редактирование объявления</FormTitle>
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
          selected={thingType() === ThingType.lost ? true : undefined}
        >
          Я потерял вещь
        </option>
        <option
          value={ThingType.found}
          selected={thingType() === ThingType.found ? true : undefined}
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
      {thingPhoto() ? (
        <ThingPhoto
          src={thingPhoto()}
          deletePhoto={() => setThingPhoto('')}
        />
      ) : (
        thingPhotoIsAvailable() && (
          <ThingPhoto
            src={pathToPhoto}
            title={`${thingName()} (изображение)`}
            deletePhoto={() => {
              deleteThingPhoto({ thingId: props.thingId }).then(() => {
                setThingPhotoIsAvailable(false);
              });
            }}
          />
        )
      )}
      <SubmitButton name='add_thing_submit'>Сохранить изменения</SubmitButton>
    </Form>
  );
};

export default UserThingEditForm;
