const menuButton = document.querySelector(".menu-toggle");
const navigation = document.querySelector("#site-nav");

if (menuButton && navigation) {
  menuButton.addEventListener("click", () => {
    const isOpen = menuButton.getAttribute("aria-expanded") === "true";
    menuButton.setAttribute("aria-expanded", String(!isOpen));
    menuButton.setAttribute("aria-label", isOpen ? "Open navigation" : "Close navigation");
    navigation.classList.toggle("is-open", !isOpen);
  });

  navigation.querySelectorAll("a").forEach((link) => {
    link.addEventListener("click", () => {
      menuButton.setAttribute("aria-expanded", "false");
      menuButton.setAttribute("aria-label", "Open navigation");
      navigation.classList.remove("is-open");
    });
  });
}

const copyButton = document.querySelector("#copy-command");
const cloneCommand = document.querySelector("#clone-command");
const copyStatus = document.querySelector("#copy-status");

if (copyButton && cloneCommand && copyStatus) {
  copyButton.addEventListener("click", async () => {
    try {
      await navigator.clipboard.writeText(cloneCommand.textContent.trim());
      copyButton.innerHTML = 'COPIED <span aria-hidden="true">✓</span>';
      copyStatus.textContent = "CLONE COMMAND COPIED TO CLIPBOARD";
    } catch {
      copyStatus.textContent = "COPY UNAVAILABLE. SELECT THE COMMAND TO COPY IT.";
    }
  });
}

const year = document.querySelector("#year");
if (year) year.textContent = new Date().getFullYear();
