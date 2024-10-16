# saving and uploading photos into superbase object storage

```
npm install @supabase/supabase-js base64-arraybuffer
```

`Replace 'YOUR_SUPABASE_URL' and 'YOUR_SUPABASE_ANON_KEY' with your actual Supabase project URL and anon key.`

This example uses Expo's ImagePicker to select an image from the device, but you can adapt it to work with other image picking libraries or your own implementation.
Remember to handle permissions for accessing the device's photo library in your React Native app.
Would you like me to explain any part of this code in more detail?
```
import { uploadImage } from './path/to/supabaseApi';
import * as ImagePicker from 'expo-image-picker';

// ... inside your component

const handleImageUpload = async () => {
  try {
    // Pick an image from the device
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ImagePicker.MediaTypeOptions.Images,
      allowsEditing: true,
      aspect: [4, 3],
      quality: 1,
      base64: true,
    });

    if (!result.canceled && result.assets[0].base64) {
      const fileName = `product_${Date.now()}.jpg`;
      const imageUrl = await uploadImage(result.assets[0].base64, fileName);
      console.log('Uploaded image URL:', imageUrl);
      // You can now use this URL to display the image or save it to your database
    }
  } catch (error) {
    console.error('Error uploading image:', error);
  }
};
```

```
import { createClient } from '@supabase/supabase-js';
import { decode } from 'base64-arraybuffer';

// Initialize Supabase client
const supabase = createClient('YOUR_SUPABASE_URL', 'YOUR_SUPABASE_ANON_KEY');

const uploadImage = async (base64Image, fileName) => {
  try {
    // Decode base64 image to ArrayBuffer
    const arrayBuffer = decode(base64Image);

    // Upload image to Supabase Storage
    const { data, error } = await supabase.storage
      .from('product-photos')
      .upload(`${fileName}`, arrayBuffer, {
        contentType: 'image/jpeg', // Adjust content type as needed
      });

    if (error) {
      throw error;
    }

    // Get public URL of the uploaded image
    const { data: { publicUrl }, error: urlError } = supabase.storage
      .from('product-photos')
      .getPublicUrl(data.path);

    if (urlError) {
      throw urlError;
    }

    return publicUrl;
  } catch (error) {
    console.error('Error uploading image:', error);
    throw error;
  }
};

export { uploadImage };
```