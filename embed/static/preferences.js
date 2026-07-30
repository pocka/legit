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

function setupTreeLayout() {
	const LOCAL_STORAGE_KEY = "legit_tree_layout";
	const VALID_LAYOUTS = ["compact"];

	function applyTreeLayout(layout) {
		if (layout) {
			localStorage.setItem(LOCAL_STORAGE_KEY, layout);
		} else {
			localStorage.removeItem(LOCAL_STORAGE_KEY);
		}

		for (const tree of document.getElementsByClassName("tree-row")) {
			tree.dataset.layout = layout;
		}
	}

	const saved = localStorage.getItem(LOCAL_STORAGE_KEY) || "";
	applyTreeLayout(saved);

	const selector = document.getElementById("tree-layout");
	if (selector) {
		selector.value = saved;

		selector.addEventListener("change", event => {
			const value = event.currentTarget.value;
			if (value && !VALID_LAYOUTS.includes(value)) {
				return;
			}

			applyTreeLayout(value);
		});
	}
}

function setupFileModeFormat() {
	const LOCAL_STORAGE_KEY = "legit_file_mode_format";
	const VALID_FORMATS = ["octal", "hidden"];
	const FILE_MODE_PATTERN = /^[d-]([r-][w-][x-]){3}$/;

	function applyFileMode(format) {
		if (format) {
			localStorage.setItem(LOCAL_STORAGE_KEY, format);
		} else {
			localStorage.removeItem(LOCAL_STORAGE_KEY);
		}

		for (const el of document.querySelectorAll("[data-file-mode]")) {
			const mode = el.dataset.fileMode;
			if (!FILE_MODE_PATTERN.test(mode || "")) {
				continue;
			}

			el.dataset.format = format;
			switch (format) {
				case "octal": {
					const n = 0
						+ (mode[1] === "r" ? 0o400 : 0)
						+ (mode[2] === "w" ? 0o200 : 0)
						+ (mode[3] === "x" ? 0o100 : 0)
						+ (mode[4] === "r" ? 0o40 : 0)
						+ (mode[5] === "w" ? 0o20 : 0)
						+ (mode[6] === "x" ? 0o10 : 0)
						+ (mode[7] === "r" ? 0o4 : 0)
						+ (mode[8] === "w" ? 0o2 : 0)
						+ (mode[9] === "x" ? 0o1 : 0);
					el.textContent = n.toString(8).padStart(3, "0");
					break;
				}
				default:
					el.textContent = mode;
					break;
			}
		}
	}

	const saved = localStorage.getItem(LOCAL_STORAGE_KEY) || "";
	applyFileMode(saved);

	const selector = document.getElementById("file-mode-format");
	if (selector) {
		selector.value = saved;

		selector.addEventListener("change", event => {
			const value = event.currentTarget.value;
			if (value && !VALID_FORMATS.includes(value)) {
				return;
			}

			applyFileMode(value);
		});
	}
}

function setupFileSizeFormat() {
	const LOCAL_STORAGE_KEY = "legit_file_size_format";
	const VALID_FORMATS = ["hidden"];

	function applyFileSize(format) {
		if (format) {
			localStorage.setItem(LOCAL_STORAGE_KEY, format);
		} else {
			localStorage.removeItem(LOCAL_STORAGE_KEY);
		}

		for (const el of document.querySelectorAll("[data-file-size]")) {
			const size = parseInt(el.dataset.fileSize, 10);
			if (!Number.isFinite(size)) {
				continue;
			}

			el.dataset.format = format;
		}
	}

	const saved = localStorage.getItem(LOCAL_STORAGE_KEY) || "";
	applyFileSize(saved);

	const selector = document.getElementById("file-size-format");
	if (selector) {
		selector.value = saved;

		selector.addEventListener("change", event => {
			const value = event.currentTarget.value;
			if (value && !VALID_FORMATS.includes(value)) {
				return;
			}

			applyFileSize(value);
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
setupTreeLayout();
setupFileModeFormat();
setupFileSizeFormat();
setupPreferencesDialog();
