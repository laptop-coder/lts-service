import { createSignal, Show, onCleanup, onMount, type Setter } from "solid-js";
import QRCode from "qrcode";

interface Props {
  text: string;
  setError: Setter<string>;
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
      props.setError(
        err instanceof Error
          ? err.message
          : "Не удалось сгенерировать QR-код",
      );
    }
  };

  const closeQR = () => setShowQR(false);

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Escape" && showQR()) {
      closeQR();
    }
  };

  onMount(() => {
    window.addEventListener("keydown", handleKeyDown);
    onCleanup(() => {
      window.removeEventListener("keydown", handleKeyDown);
    });
  });

  return (
    <>
      <button
        onClick={generateQR}
        class="ml-2 px-3 py-1 text-sm bg-green-700 text-white rounded hover:bg-green-800 transition cursor-pointer"
        title="Показать QR-код"
      >
        QR
      </button>

      <Show when={showQR()}>
        <div
          class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
          onClick={closeQR}
        >
          <div
            class="bg-white p-6 rounded-lg shadow-xl"
            onClick={(e) => e.stopPropagation()}
          >
            <img src={qrDataUrl()} alt="QR Code" class="w-64 h-64" />
            <div class="mt-4 text-center">
              <button
                onClick={closeQR}
                class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition cursor-pointer"
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
