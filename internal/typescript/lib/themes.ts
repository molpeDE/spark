import { useSignal } from "@preact/signals";
import { h } from 'preact'

export const preference = () => localStorage.getItem('theme') || ((window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) ? 'dark' : 'light');

export function setPreference(theme: string) {
	localStorage.setItem('theme', theme);
	document.body.setAttribute('data-bs-theme', theme);
}

export function ThemeToggle() {

	const dark = useSignal(preference() === 'dark');

	function toggle() {
		setPreference(dark.value ? 'light' : 'dark');
		dark.value = preference() === 'dark';
	}

	return h('i', { class: `clickable p-2 fa-solid fa-${dark.value ? 'moon' : 'sun'}`, onClick: toggle, title: 'Toggle Theme' }, null);
}

setPreference(preference());