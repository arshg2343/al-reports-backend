import React, { useState, useEffect } from 'react';

const GlitchReportForm = () => {
  const [formData, setFormData] = useState({
    email: '',
    username: '',
    deviceType: '',
    browserInfo: '',
    glitchType: 'visual',
    glitchLocation: '',
    glitchDescription: '',
    stepsToReproduce: '',
    attachment: null,
    urgency: 'medium'
  });
  
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [glitchEffect, setGlitchEffect] = useState(false);
  const [filePreview, setFilePreview] = useState(null);

  // Subtle glitch effect animation - less frequent for professional look
  useEffect(() => {
    const interval = setInterval(() => {
      setGlitchEffect(true);
      setTimeout(() => setGlitchEffect(false), 150);
    }, 5000);
    
    return () => clearInterval(interval);
  }, []);
  
  // More subtle binary background effect
  const [noise, setNoise] = useState('');
  useEffect(() => {
    const generateNoise = () => {
      let result = '';
      const characters = '01';
      for (let i = 0; i < 800; i++) {
        result += characters.charAt(Math.floor(Math.random() * 2));
      }
      return result;
    };
    
    const interval = setInterval(() => {
      setNoise(generateNoise());
    }, 300); // Slower update for less distraction
    
    return () => clearInterval(interval);
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setFormData(prev => ({ ...prev, attachment: file }));
      
      // Create preview for image
      if (file.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onloadend = () => {
          setFilePreview(reader.result);
        };
        reader.readAsDataURL(file);
      } else {
        setFilePreview(null);
      }
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    // Simulating standard form submission
    try {
      // In a real application, you would use FormData to handle file uploads
      const formDataToSend = new FormData();
      
      // Append all form fields
      Object.keys(formData).forEach(key => {
        if (key === 'attachment' && formData[key]) {
          formDataToSend.append(key, formData[key]);
        } else if (key !== 'attachment') {
          formDataToSend.append(key, formData[key]);
        }
      });
      
      // Simulate API call with timeout
      setTimeout(() => {
        console.log('Form data would be submitted:', formDataToSend);
        setIsSubmitted(true);
        setIsSubmitting(false);
      }, 1500);
      
    } catch (error) {
      console.error('Error submitting form:', error);
      setIsSubmitting(false);
    }
  };

  const clearForm = () => {
    setFormData({
      email: '',
      username: '',
      deviceType: '',
      browserInfo: '',
      glitchType: 'visual',
      glitchLocation: '',
      glitchDescription: '',
      stepsToReproduce: '',
      attachment: null,
      urgency: 'medium'
    });
    setFilePreview(null);
    setIsSubmitted(false);
  };

  return (
    <div className="min-h-screen bg-gray-900 text-blue-400 py-12 relative overflow-hidden font-sans">
      {/* Subtle Binary Noise Background */}
      <div className="fixed inset-0 opacity-3 text-xs z-0 overflow-hidden whitespace-pre-wrap text-blue-500/30">
        {noise}
      </div>
      
      {/* Main Content */}
      <div className="container mx-auto px-4 relative z-10">
        <header className="mb-12 text-center">
          <h1 
            className={`text-4xl md:text-5xl font-bold mb-4 text-blue-300 tracking-wide ${glitchEffect ? 'relative overflow-hidden' : ''}`}
            style={glitchEffect ? {
              textShadow: '1px 0 #6366f1, -1px 0 #2dd4bf',
              animation: 'shake 0.15s linear'
            } : {}}
          >
            SYSTEM ISSUE REPORT
          </h1>
          <div className="w-24 h-1 bg-blue-500 mx-auto mb-6"></div>
          <p className="text-lg max-w-2xl mx-auto text-gray-300 font-light">
            Please provide detailed information about the system anomaly you've encountered.
            All reports are prioritized according to severity and impact on operations.
          </p>
        </header>

        {isSubmitted ? (
          <div className="max-w-3xl mx-auto bg-gray-800 p-8 rounded-lg border border-blue-400 shadow-lg shadow-blue-500/10">
            <div className="text-center">
              <div className="inline-block rounded-full bg-blue-900/50 p-4 mb-6">
                <svg className="w-16 h-16 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7"></path>
                </svg>
              </div>
              <h2 className="text-3xl font-bold text-blue-300 mb-4">Report Submitted Successfully</h2>
              <p className="text-xl mb-6 text-gray-300">Your issue has been logged and will be addressed promptly.</p>
              <p className="text-blue-200 mb-8">Reference ID: {Math.random().toString(36).substring(2, 12).toUpperCase()}</p>
              <div className="flex justify-center gap-4 flex-wrap">
                <button 
                  onClick={clearForm}
                  className="bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-6 rounded-md transition duration-300"
                >
                  Submit Another Report
                </button>
                <button 
                  onClick={() => window.print()}
                  className="bg-gray-700 hover:bg-gray-600 text-white font-medium py-3 px-6 rounded-md transition duration-300"
                >
                  Print Confirmation
                </button>
              </div>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="max-w-3xl mx-auto bg-gray-800 p-8 rounded-lg border border-blue-400 shadow-lg shadow-blue-500/10">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
              <div>
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="email">Email Address</label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  required
                  value={formData.email}
                  onChange={handleChange}
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                  placeholder="user@company.com"
                />
              </div>
              
              <div>
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="username">Username/ID</label>
                <input
                  type="text"
                  id="username"
                  name="username"
                  required
                  value={formData.username}
                  onChange={handleChange}
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                  placeholder="Your system username"
                />
              </div>

              <div>
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="deviceType">Device Type</label>
                <input
                  type="text"
                  id="deviceType"
                  name="deviceType"
                  required
                  value={formData.deviceType}
                  onChange={handleChange}
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                  placeholder="e.g. iPhone 13, Dell XPS 15"
                />
              </div>

              <div>
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="browserInfo">Browser Info</label>
                <input
                  type="text"
                  id="browserInfo"
                  name="browserInfo"
                  required
                  value={formData.browserInfo}
                  onChange={handleChange}
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                  placeholder="e.g. Chrome 99, Firefox 95"
                />
              </div>

              <div>
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="glitchType">Issue Type</label>
                <select
                  id="glitchType"
                  name="glitchType"
                  required
                  value={formData.glitchType}
                  onChange={handleChange}
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                >
                  <option value="visual">Visual Issue</option>
                  <option value="functional">Functional Error</option>
                  <option value="performance">Performance Issue</option>
                  <option value="crash">System Crash</option>
                  <option value="security">Security Vulnerability</option>
                  <option value="data">Data Inconsistency</option>
                  <option value="other">Other</option>
                </select>
              </div>

              <div>
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="glitchLocation">Issue Location</label>
                <input
                  type="text"
                  id="glitchLocation"
                  name="glitchLocation"
                  required
                  value={formData.glitchLocation}
                  onChange={handleChange}
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                  placeholder="e.g. Login Page, Dashboard"
                />
              </div>

              <div className="md:col-span-2">
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="glitchDescription">Detailed Description</label>
                <textarea
                  id="glitchDescription"
                  name="glitchDescription"
                  required
                  value={formData.glitchDescription}
                  onChange={handleChange}
                  rows="4"
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                  placeholder="Provide detailed information about what happened..."
                ></textarea>
              </div>

              <div className="md:col-span-2">
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="stepsToReproduce">Steps to Reproduce</label>
                <textarea
                  id="stepsToReproduce"
                  name="stepsToReproduce"
                  required
                  value={formData.stepsToReproduce}
                  onChange={handleChange}
                  rows="4"
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                  placeholder="1. Logged in as user\n2. Navigated to...\n3. Clicked on..."
                ></textarea>
              </div>

              <div className="md:col-span-2">
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="fileUpload">
                  Attach Screenshot/Video
                </label>
                <div className="flex items-center justify-center w-full">
                  <label className="flex flex-col w-full h-32 border-2 border-dashed border-blue-500/50 rounded-lg cursor-pointer hover:bg-gray-800/50 transition-all">
                    <div className="flex flex-col items-center justify-center pt-7">
                      {filePreview ? (
                        <img src={filePreview} alt="Preview" className="h-16 object-contain" />
                      ) : (
                        <>
                          <svg className="w-8 h-8 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
                          </svg>
                          <p className="pt-1 text-sm tracking-wider text-gray-300">
                            Drag & drop or click to attach file
                          </p>
                          <p className="text-xs text-gray-400">
                            (PNG, JPG, GIF, MP4 up to 10MB)
                          </p>
                        </>
                      )}
                    </div>
                    <input
                      type="file"
                      id="fileUpload"
                      name="attachment"
                      onChange={handleFileChange}
                      className="opacity-0"
                      accept="image/*, video/*"
                    />
                  </label>
                </div>
                {formData.attachment && (
                  <p className="mt-2 text-sm text-blue-300">
                    File selected: {formData.attachment.name}
                  </p>
                )}
              </div>

              <div className="md:col-span-2">
                <label className="block text-blue-300 mb-2 font-medium" htmlFor="urgency">Priority Level</label>
                <select
                  id="urgency"
                  name="urgency"
                  required
                  value={formData.urgency}
                  onChange={handleChange}
                  className="w-full bg-gray-900 border border-blue-500/50 rounded-md py-3 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                >
                  <option value="low">Low - Minor Issue (No Impact on Operation)</option>
                  <option value="medium">Medium - Moderate Impact on Usability</option>
                  <option value="high">High - Severe Impact on Critical Functionality</option>
                  <option value="critical">Critical - Complete System Failure</option>
                </select>
              </div>
            </div>

            <div className="text-center">
              <button
                type="submit"
                disabled={isSubmitting}
                className={`bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-8 rounded-md transition duration-300 ${isSubmitting ? 'opacity-50 cursor-not-allowed' : ''}`}
              >
                {isSubmitting ? (
                  <span className="flex items-center justify-center">
                    <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Processing...
                  </span>
                ) : (
                  'Submit Report'
                )}
              </button>
            </div>
          </form>
        )}

        <div className="mt-12 text-center text-gray-400 text-sm">
          <p>System Monitoring Platform v3.2.1</p>
          <p>All reports are encrypted and securely stored for analysis.</p>
        </div>
      </div>

      {/* Subtle decorative elements */}
      <div className="fixed top-0 left-0 w-full h-0.5 bg-blue-500 opacity-60"></div>
      <div className="fixed bottom-0 left-0 w-full h-0.5 bg-blue-500 opacity-60"></div>
      <div className="fixed top-0 left-0 h-full w-0.5 bg-blue-500 opacity-60"></div>
      <div className="fixed top-0 right-0 h-full w-0.5 bg-blue-500 opacity-60"></div>
    </div>
  );
};

export default GlitchReportForm;