document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('.scanner-form');

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const requestData = {
            host: document.getElementById('host').value,
            ports: document.getElementById('ports').value,
            protocol: document.getElementById('protocol').value,
            timeout: document.getElementById('timeout').value
        };

        try {
            const response = await fetch('/scan', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestData)
            });

            if (!response.ok) {
                throw new Error('Scan failed');
            }

            const html = await response.text();
            document.open();
            document.write(html);
            document.close();
        } catch (error) {
            alert('Error: ' + error.message);
            console.error('Scan error:', error);
        }
    });
});