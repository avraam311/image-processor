document.getElementById('uploadForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const formData = new FormData();
    const imageFile = document.getElementById('image').files[0];
    const processing = document.getElementById('processing').value;

    if (!imageFile) {
        showStatus('Please select an image file.', 'error');
        return;
    }

    formData.append('image', imageFile);
    formData.append('processing', processing);

    showStatus('Uploading and processing image...', '');

    try {
        // Upload image
        const uploadResponse = await fetch('http://localhost:8080/image-processor/api/upload', {
            method: 'POST',
            body: formData
        });

        if (!uploadResponse.ok) {
            throw new Error(`Upload failed: ${uploadResponse.statusText}`);
        }

        const uploadResult = await uploadResponse.json();
        const imageId = uploadResult.id;

        showStatus(`Image uploaded with ID: ${imageId}. Processing...`, 'success');

        // Poll for processed image
        pollForImage(imageId);

    } catch (error) {
        showStatus(`Error: ${error.message}`, 'error');
    }
});

function pollForImage(id) {
    const pollInterval = setInterval(async () => {
        try {
            const response = await fetch(`http://localhost:8080/image-processor/api/image/${id}`);

            if (response.status === 200) {
                clearInterval(pollInterval);
                const blob = await response.blob();
                const imageUrl = URL.createObjectURL(blob);
                document.getElementById('result').innerHTML = `<h2>Processed Image:</h2><img src="${imageUrl}" alt="Processed Image">`;
                showStatus('Image processed successfully!', 'success');
            } else if (response.status === 503) {
                showStatus('Image is still processing...', '');
            } else {
                throw new Error(`Failed to get image: ${response.statusText}`);
            }
        } catch (error) {
            clearInterval(pollInterval);
            showStatus(`Error: ${error.message}`, 'error');
        }
    }, 2000); // Poll every 2 seconds
}

function showStatus(message, className) {
    const statusDiv = document.getElementById('status');
    statusDiv.textContent = message;
    statusDiv.className = className;
}
