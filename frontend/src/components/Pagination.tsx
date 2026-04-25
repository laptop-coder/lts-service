interface PaginationProps {
  page: number;
  hasMore: boolean;
  onPrev: () => void;
  onNext: () => void;
}

const Pagination = (props: PaginationProps) => {
  return (
    <div class="flex justify-center items-center gap-4 mt-6">
      <button
        onClick={props.onPrev}
        disabled={props.page === 0}
        class="px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
      >
        ← Назад
      </button>
      <span class="text-sm text-gray-500">Страница {props.page + 1}</span>
      <button
        onClick={props.onNext}
        disabled={!props.hasMore}
        class="px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
      >
        Вперёд →
      </button>
    </div>
  );
};

export default Pagination
