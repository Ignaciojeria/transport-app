/**
 * Utility functions for uploading photos to delivery unit evidence URLs
 */

export interface DeliveryUnit {
  evidences?: Array<{
    downloadUrl?: string;
    uploadUrl?: string;
  }>;
}

export interface RouteData {
  visits?: Array<{
    orders?: Array<{
      deliveryUnits?: DeliveryUnit[];
    }>;
  }>;
}

/**
 * Converts an image to WebP format
 * @param imageBlob The original image blob
 * @param quality Quality of the WebP image (0.1 to 1.0)
 * @returns Promise<Blob> WebP image blob
 */
async function convertToWebP(imageBlob: Blob, quality: number = 0.8): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    const img = new Image();
    
    img.onload = () => {
      canvas.width = img.width;
      canvas.height = img.height;
      
      ctx?.drawImage(img, 0, 0);
      
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve(blob);
          } else {
            reject(new Error('Failed to convert image to WebP'));
          }
        },
        'image/webp',
        quality
      );
    };
    
    img.onerror = () => reject(new Error('Failed to load image'));
    img.src = URL.createObjectURL(imageBlob);
  });
}

/**
 * Uploads a photo to the delivery unit's evidence upload URL
 * @param photoDataUrl Base64 data URL of the photo
 * @param uploadUrl The upload URL from the delivery unit evidence
 * @returns Promise<boolean> Success status
 */
export async function uploadPhotoToEvidence(
  photoDataUrl: string, 
  uploadUrl: string
): Promise<boolean> {
  try {
    // Convert base64 data URL to blob
    const response = await fetch(photoDataUrl);
    const originalBlob = await response.blob();
    
    // Convert to WebP format
    console.log('🔄 Converting image to WebP format...');
    const webpBlob = await convertToWebP(originalBlob, 0.8);
    console.log('✅ Image converted to WebP successfully');
    
    // Upload the WebP photo to the provided URL
    const uploadResponse = await fetch(uploadUrl, {
      method: 'PUT',
      body: webpBlob,
      headers: {
        'Content-Type': 'image/webp',
      },
    });

    if (uploadResponse.ok) {
      console.log('✅ WebP photo uploaded successfully to:', uploadUrl);
      return true;
    } else {
      console.error('❌ Failed to upload WebP photo:', uploadResponse.status, uploadResponse.statusText);
      return false;
    }
  } catch (error) {
    console.error('❌ Error uploading photo:', error);
    return false;
  }
}

/**
 * Gets the upload URL for a specific delivery unit
 * @param routeData The route data containing visits
 * @param visitIndex Visit index
 * @param orderIndex Order index
 * @param unitIndex Delivery unit index
 * @returns Upload URL or undefined if not found
 */
export function getDeliveryUnitUploadUrl(
  routeData: RouteData,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): string | undefined {
  const visit = routeData?.visits?.[visitIndex];
  const order = visit?.orders?.[orderIndex];
  const deliveryUnit = order?.deliveryUnits?.[unitIndex];
  const evidence = deliveryUnit?.evidences?.[0]; // Assuming first evidence entry
  
  return evidence?.uploadUrl;
}

/**
 * Checks if a Storj URL is valid (not expired)
 * @param url The Storj URL to validate
 * @returns boolean indicating if URL is valid
 */
function isStorjUrlValid(url: string): boolean {
  if (!url) return false;
  
  try {
    const urlObj = new URL(url);
    const expiresParam = urlObj.searchParams.get('X-Amz-Expires');
    
    if (!expiresParam) return false;
    
    const expiresSeconds = parseInt(expiresParam);
    // Check if expires is less than a week (604800 seconds)
    if (expiresSeconds >= 604800) return false;
    
    // Check if the URL hasn't expired yet
    const dateParam = urlObj.searchParams.get('X-Amz-Date');
    if (dateParam) {
      const urlDate = new Date(
        dateParam.replace(/(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})(\d{2})Z/, '$1-$2-$3T$4:$5:$6Z')
      );
      const expirationTime = urlDate.getTime() + (expiresSeconds * 1000);
      if (Date.now() > expirationTime) return false;
    }
    
    return true;
  } catch (error) {
    console.warn('Error validating Storj URL:', error);
    return false;
  }
}

/**
 * Gets the download URL for a specific delivery unit
 * @param routeData The route data containing visits
 * @param visitIndex Visit index
 * @param orderIndex Order index
 * @param unitIndex Delivery unit index
 * @returns Download URL or undefined if not found or invalid
 */
export function getDeliveryUnitDownloadUrl(
  routeData: RouteData,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): string | undefined {
  const visit = routeData?.visits?.[visitIndex];
  const order = visit?.orders?.[orderIndex];
  const deliveryUnit = order?.deliveryUnits?.[unitIndex];
  const evidence = deliveryUnit?.evidences?.[0]; // Assuming first evidence entry
  
  const downloadUrl = evidence?.downloadUrl;
  
  // Validate the URL before returning it
  if (downloadUrl && isStorjUrlValid(downloadUrl)) {
    return downloadUrl;
  }
  
  // Return a note indicating the URL is invalid/expired
  if (downloadUrl) {
    console.warn('Storj download URL is invalid or expired:', downloadUrl);
    return 'URL_EXPIRADA_O_INVALIDA';
  }
  
  return undefined;
}