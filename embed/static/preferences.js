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
setupPreferencesDialog();
