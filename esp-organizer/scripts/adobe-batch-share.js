const puppeteer = require('puppeteer');
const fs = require('fs');

async function createAdobeShareLinks() {
    const browser = await puppeteer.launch({ 
        headless: false, // Show browser for debugging
        slowMo: 100 // Slow down for stability
    });
    
    const page = await browser.newPage();
    
    // Login to Adobe Document Cloud
    await page.goto('https://documentcloud.adobe.com/');
    
    console.log('⏳ Please log in manually...');
    await page.waitForNavigation({ timeout: 120000 }); // Wait 2 minutes for manual login
    
    console.log('✅ Logged in! Starting batch link creation...');
    
    const shareLinks = [];
    
    // Get all PDF files in the interface
    const pdfElements = await page.$$('[data-test-id="document-item"]');
    
    console.log(`📄 Found ${pdfElements.length} documents`);
    
    for (let i = 0; i < pdfElements.length; i++) {
        try {
            // Click on document to select it
            await pdfElements[i].click();
            await page.waitForTimeout(500);
            
            // Click "Share" button
            const shareButton = await page.$('[data-test-id="share-button"]');
            if (shareButton) {
                await shareButton.click();
                await page.waitForTimeout(1000);
                
                // Click "Create Link"
                const createLinkButton = await page.$('[data-test-id="create-link"]');
                if (createLinkButton) {
                    await createLinkButton.click();
                    await page.waitForTimeout(1000);
                    
                    // Copy the link
                    const linkInput = await page.$('[data-test-id="share-link-input"]');
                    if (linkInput) {
                        const link = await page.evaluate(el => el.value, linkInput);
                        shareLinks.push(link);
                        console.log(`✅ Created link ${i + 1}: ${link}`);
                    }
                    
                    // Close the share dialog
                    const closeButton = await page.$('[data-test-id="close-dialog"]');
                    if (closeButton) await closeButton.click();
                }
            }
            
            await page.waitForTimeout(1000);
        } catch (error) {
            console.error(`❌ Error processing document ${i + 1}:`, error.message);
        }
    }
    
    // Save links to file
    const linksContent = shareLinks.map((link, i) => `# Case ${i + 1}\n${link}`).join('\n\n');
    fs.writeFileSync('resources/immunology/adobe_scan_links.txt', linksContent);
    
    console.log(`\n✅ Created ${shareLinks.length} share links`);
    console.log(`📄 Saved to: resources/immunology/adobe_scan_links.txt`);
    
    await browser.close();
}

createAdobeShareLinks();
