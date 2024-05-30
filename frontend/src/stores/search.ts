import { create } from 'zustand';

interface SearchState {
  search: string;
  setSearch(text: string): void;
}

const searchStore = create<SearchState>((set) => ({
  search: '',
  setSearch(text: string): void {
    set({ search: text });
  },
}));

export default searchStore;
