import { createSignal, Show } from "solid-js";
import QRCode from "qrcode";

interface Props {
  text: string;
}

const QRCodeButton = (props: Props) => {
  const [showQR, setShowQR] = createSignal(false);
  const [qrDataUrl, setQrDataUrl] = createSignal("");

  const generateQR = async () => {
    if (qrDataUrl()) {
      setShowQR(!showQR());
      return;
    }
    try {
      const url = await QRCode.toDataURL(props.text);
      setQrDataUrl(url);
      setShowQR(true);
    } catch (err) {
      // TODO
      console.error("Failed to generate QR code");
    }
  };

  const closeQR = () => setShowQR(false);

  return (
    <>
      <button
        onClick={generateQR}
        class="ml-2 px-3 py-1 text-sm bg-green-500 text-white rounded hover:bg-green-600"
        title="Показать QR-код"
      >
        QR
      </button>

      <Show when={showQR()}>
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={closeQR}>
          <div class="bg-white p-6 rounded-lg shadow-xl" onClick={(e) => e.stopPropagation()}>
            <img src={qrDataUrl()} alt="QR Code" class="w-64 h-64" />
            <div class="mt-4 text-center">
              <button
                onClick={closeQR}
                class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
              >
                Закрыть
              </button>
            </div>
          </div>
        </div>
      </Show>
    </>
  );
};

export default QRCodeButton;
