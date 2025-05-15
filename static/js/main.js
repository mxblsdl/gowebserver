async function downloadFile(file_id) {
  try {
    const apiKey = document.querySelector("#key").value;

    const response = await fetch(`/download/${file_id}`, {
      method: "GET",
      headers: { access_token: apiKey },
    });

    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const disposition = response.headers.get("Content-Disposition");
    const filename = disposition.match(/filename="?([^"]*)"?/)[1];

    // Create a Blob from the response
    const blob = await response.blob();

    // Create a link element, set the download attribute and click it
    downloadBlob(blob, filename);
  } catch (error) {
    console.error("Error:", error);
    alertify.error("Failed to download file.");
  }
}

function downloadBlob(blob, filename) {
  const link = document.createElement("a");
  const objectURL = URL.createObjectURL(blob);
  link.href = objectURL;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(objectURL); // Clean up
}

function copyToClipboard(input) {
  let textToCopy;

  if (typeof input === "string") {
    textToCopy = input;
  } else if (input instanceof Event) {
    const response = input.detail.xhr.response;
    const data = JSON.parse(response);
    textToCopy = data.link;
  } else {
    console.error("Invalid input to copyToClipboard");
    return;
  }

  navigator.clipboard
    .writeText(textToCopy)
    .then(() => {
      alertify.success("Link copied to clipboard!");
    })
    .catch((error) => {
      console.error("Error copying to clipboard:", error);
    });
}
