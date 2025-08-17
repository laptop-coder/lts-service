import { JSX, createSignal } from 'solid-js';

import Header from '../components/Header/Header';
import Content from '../components/Content/Content';
import Footer from '../components/Footer/Footer';
import Page from '../ui/Page/Page';
import SquareImageButton from '../ui/SquareImageButton/SquareImageButton';
import SubmitButton from '../ui/SubmitButton/SubmitButton';
import ResetButton from '../ui/ResetButton/ResetButton';
import Select from '../ui/Select/Select';
import Input from '../ui/Input/Input';
import CenterForm from '../ui/CenterForm/CenterForm';
import TextArea from '../ui/TextArea/TextArea';
import AttachFile from '../ui/AttachFile/AttachFile';
import FormButtonsBlock from '../ui/FormButtonsBlock/FormButtonsBlock';
import ThingPhoto from '../ui/ThingPhoto/ThingPhoto';
import fileToBase64 from '../utils/fileToBase64';
import axiosInstanceUnauthorized from '../utils/axiosInstanceUnauthorized';

const AddThingPage = (): JSX.Element => {
  const [thingType, setThingType] = createSignal('lost');
  const [thingName, setThingName] = createSignal('');
  const [userEmail, setUserEmail] = createSignal('');
  const [thingLocation, setThingLocation] = createSignal('');
  const [customText, setCustomText] = createSignal('');
  const [thingPhoto, setThingPhoto] = createSignal('');

  // TODO: refactor
  const fieldsAreNotFilledMessage = () =>
    alert('Обязательные поля не заполнены (они отмечены звёздочкой*)');
  const sendingErrorMessage = () =>
    alert('Ошибка отправки. Попробуйте ещё раз');
  const wrongThingTypeMessage = () => console.log('Ошибка. Неверный тип вещи');

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    const [data, setData] = createSignal({});
    if (thingType() === 'lost') {
      if (thingName() !== '' && userEmail() !== '') {
        setData({
          thingName: thingName(),
          userEmail: userEmail(),
          customText: customText(),
        });
      } else {
        fieldsAreNotFilledMessage();
        return;
      }
    } else if (thingType() === 'found') {
      if (thingName() !== '' && thingLocation() !== '') {
        setData({
          thingName: thingName(),
          thingLocation: thingLocation(),
          customText: customText(),
        });
      } else {
        fieldsAreNotFilledMessage();
        return;
      }
    } else {
      wrongThingTypeMessage();
      return;
    }
    await axiosInstanceUnauthorized
      .post(`/thing/add?thing_type=${thingType()}`, data(), {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      })
      .then((response) => {
        if (response.status !== 200) {
          sendingErrorMessage();
        } else {
          window.location.replace(
            `/thing_status?thing_type=${thingType()}&thing_id=${response.data.ThingId}`,
          );
        }
      })
      .catch((error) => console.log(error));
  };

  return (
    <Page>
      <Header>
        {window.history.length > 1 && (
          <SquareImageButton onclick={() => window.history.back()}>
            <img src='/src/assets/arrow_back.svg' />
          </SquareImageButton>
        )}
      </Header>
      <Content>
        <CenterForm
          header='Добавление новой вещи'
          method='post'
          onsubmit={(event) => handleSubmit(event)}
        >
          <Select
            id='thingTypeSelect'
            value={thingType()}
            onchange={(event) => setThingType(event.target.value)}
            label='Что произошло?*'
          >
            <option value='lost'>Я потерял свою вещь</option>
            <option value='found'>Я нашёл чью-то вещь</option>
          </Select>
          {thingType() === 'lost' && (
            <>
              <Input
                placeholder='Что Вы потеряли?*'
                value={thingName()}
                onchange={(event) => setThingName(event.target.value)}
              />
              <Input
                type='email'
                placeholder='Email для связи с Вами*'
                value={userEmail()}
                onchange={(event) => setUserEmail(event.target.value)}
              />
            </>
          )}
          {thingType() === 'found' && (
            <>
              <Input
                placeholder='Что Вы нашли?*'
                value={thingName()}
                onchange={(event) => setThingName(event.target.value)}
              />
              <Input
                placeholder='Где забрать вещь?*'
                value={thingLocation()}
                onchange={(event) => setThingLocation(event.target.value)}
              />
            </>
          )}
          <TextArea
            placeholder='Здесь можно оставить сообщение'
            value={customText()}
            onchange={(event) => setCustomText(event.target.value)}
          />
          <AttachFile
            accept='image/jpeg,image/png'
            id='attachThingPhoto'
            label='Выберите фото'
            onchange={(event) =>
              event.target.files &&
              fileToBase64(event.target.files[0]).then((photoBase64) =>
                setThingPhoto(photoBase64),
              )
            }
          />
          <ThingPhoto src={thingPhoto()} />
          <FormButtonsBlock>
            <ResetButton onclick={() => setThingPhoto('')}>
              Сбросить
            </ResetButton>
            <SubmitButton>Отправить</SubmitButton>
          </FormButtonsBlock>
        </CenterForm>
      </Content>
      <Footer />
    </Page>
  );
};

export default AddThingPage;
