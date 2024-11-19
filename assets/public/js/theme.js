(() => {
  const themeStitcher = document.getElementById("themingSwitcher");

  // set toggler position based on system theme
  if (isSystemThemeSetToDark) {
    themeStitcher.checked = true;
  }

  document.getElementById("darkModeItem").addEventListener("click", (e) => {
    const themeStitcher = document.getElementById("themingSwitcher");

    themeStitcher.checked = !themeStitcher.checked;
    toggleTheme(themeStitcher.checked);
  });

  const toggleTheme = (isChecked) => {
    const theme = isChecked ? "dark" : "light";

    document.documentElement.dataset.mdbTheme = theme;
  };

  // add listener to toggle theme with Shift + D
  document.addEventListener("keydown", (e) => {
    if (e.shiftKey && e.key === "D") {
      themeStitcher.checked = !themeStitcher.checked;
      toggleTheme(themeStitcher.checked);
    }
  });
})();
