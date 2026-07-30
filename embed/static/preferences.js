// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

function setupTabWidthSelector() {
	const LOCAL_STORAGE_KEY = "legit_tab_width";

	function isValidTabWidth(width) {
		if (typeof width !== "string") {
			return false;
		}

		const n = parseInt(width, 10);

		return Number.isFinite(n) && n > 0;
	}

	function applyTabWidth(width) {
		localStorage.setItem(LOCAL_STORAGE_KEY, width);
		document.body.style.tabSize = width;
	}

	const saved = localStorage.getItem(LOCAL_STORAGE_KEY);
	const value = isValidTabWidth(saved) ? saved : "4";
	applyTabWidth(value);

	const selector = document.getElementById("global-tab-width");
	if (selector) {
		selector.value = value;

		selector.addEventListener("change", event => {
			if (!isValidTabWidth(event.currentTarget.value)) {
				return;
			}

			applyTabWidth(event.currentTarget.value);
		});
	}
}

function setupThemeColor() {
	const LOCAL_STORAGE_KEY = "legit_theme_color";
	const VALID_COLORS = ["red", "orange", "green", "gray", "pink", "purple", "yellow"];

	function applyThemeColor(color) {
		if (color) {
			localStorage.setItem(LOCAL_STORAGE_KEY, color);
			document.documentElement.dataset.themeColor = color;
		} else {
			localStorage.removeItem(LOCAL_STORAGE_KEY);
			document.documentElement.dataset.themeColor = "";
		}
	}

	const saved = localStorage.getItem(LOCAL_STORAGE_KEY) || "";
	applyThemeColor(saved);

	const selector = document.getElementById("theme-color");
	if (selector) {
		selector.value = saved;

		selector.addEventListener("change", event => {
			const value = event.currentTarget.value;
			if (value && !VALID_COLORS.includes(value)) {
				return;
			}

			applyThemeColor(value);
		});
	}
}

function setupPreferencesDialog() {
	const preferences = document.getElementById("preferences-dialog");
	document.getElementById("preferences-dialog-trigger")?.addEventListener("click", (event) => {
		preferences?.showModal();
	});
	document.getElementById("preferences-dialog-closer")?.addEventListener("click", (event) => {
		preferences?.close();
	});
}

setupTabWidthSelector();
setupThemeColor();
setupPreferencesDialog();
