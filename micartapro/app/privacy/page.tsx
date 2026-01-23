import Link from 'next/link'
import Image from 'next/image'

export const metadata = {
  title: 'Privacy Policy - MiCartaPro',
  description: 'MiCartaPro privacy policy. Learn how we protect and manage your personal data.',
}

export default function PrivacyPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
      {/* Navigation */}
      <nav className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link href="/" className="flex items-center space-x-2">
              <Image 
                src="/logo.png" 
                alt="MiCartaPro Logo" 
                width={240} 
                height={72}
                className="h-16 md:h-20 w-auto"
              />
            </Link>
            <Link 
              href="/"
              className="text-gray-600 hover:text-blue-600 transition-colors"
            >
              Back to Home
            </Link>
          </div>
        </div>
      </nav>

      {/* Content */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-16">
        <div className="bg-white rounded-lg shadow-lg p-8 md:p-12">
          <h1 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            Privacy Policy
          </h1>
          
          <p className="text-gray-600 mb-8">
            <strong>Last updated:</strong> {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>

          <div className="prose prose-lg max-w-none space-y-8 text-gray-700">
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. Introduction</h2>
              <p>
                MiCartaPro ("we", "our", or "the company") is committed to protecting and respecting your privacy. This Privacy Policy explains how we collect, use, disclose, and protect your personal information when you use our digital menu services.
              </p>
              <p>
                <strong>Data Controller:</strong> MiCartaPro (micartapro.com) is the data controller responsible for processing your personal information in accordance with this Privacy Policy.
              </p>
              <p>
                By using our services, you agree to the practices described in this policy. If you do not agree with this policy, please do not use our services.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Information We Collect</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">2.1. Information You Provide Directly</h3>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Order information:</strong> Customer name, pickup time, and complete order details when sent through WhatsApp.</li>
                <li><strong>Contact information:</strong> Phone number when you contact us through WhatsApp for quotes or inquiries.</li>
                <li><strong>Restaurant information:</strong> Data provided when requesting our services, including contact information and business details.</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">2.2. Automatically Collected Information</h3>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Cart information:</strong> Selected items are stored locally in your browser (localStorage) to improve your user experience.</li>
                <li><strong>Device information:</strong> Device type, operating system, and browser used to access our services.</li>
                <li><strong>Usage information:</strong> How you interact with our services, including pages visited and time spent.</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Legal Basis for Processing (GDPR/LGPD)</h2>
              <p>We process your personal information based on the following legal bases:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Consent:</strong> When you explicitly provide your consent for the processing of your data.</li>
                <li><strong>Contract performance:</strong> To process and manage orders from restaurants using our services.</li>
                <li><strong>Legitimate interest:</strong> To improve our services, prevent fraud, and maintain the security of our systems.</li>
                <li><strong>Legal compliance:</strong> When necessary to comply with applicable legal obligations.</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Use of Information</h2>
              <p>We use the collected information for the following purposes:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Process and manage restaurant orders efficiently</li>
                <li>Improve and personalize user experience</li>
                <li>Respond to inquiries, quote requests, and provide customer support</li>
                <li>Maintain shopping cart functionality and remember your preferences</li>
                <li>Develop and improve our services and features</li>
                <li>Send communications related to our services (with your consent)</li>
                <li>Comply with legal obligations and protect our legal rights</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Data Storage and Retention</h2>
              <p>
                <strong>Local Storage:</strong> Cart data is stored locally in your browser (localStorage) and is not sent to our servers until you decide to send an order through WhatsApp. You can clear this data at any time from your browser settings.
              </p>
              <p>
                <strong>Data Retention:</strong> We retain your personal information only for as long as necessary to fulfill the purposes described in this policy, unless the law requires or permits a longer retention period.
              </p>
              <p>
                <strong>Order Data:</strong> Orders are processed directly through WhatsApp and are managed by the corresponding restaurant. We do not store order information on our servers.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. Sharing and Disclosure of Information</h2>
              <p>
                <strong>We never sell personal data to third parties.</strong> We do not sell, rent, or share your personal information with third parties for their own commercial purposes, except in the following circumstances:
              </p>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.1. Service Providers</h3>
              <p>We may share information with trusted service providers who help us operate our business, including:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Hosting and cloud service providers (e.g., Firebase, Vercel, or similar infrastructure providers)</li>
                <li>Payment processors (who act as Merchant of Record)</li>
                <li>Communication service providers (WhatsApp)</li>
              </ul>
              <p className="mt-4">
                All service providers are contractually obligated to protect your information and use it only for the purposes we specify.
              </p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.2. Restaurants</h3>
              <p>When you send an order, the information is shared with the corresponding restaurant through WhatsApp to process your order.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.3. Legal Requirements</h3>
              <p>We may disclose information when required by law, court order, or to protect our legal rights, property, or safety, or that of our users.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.4. Business Transfers</h3>
              <p>In the event of a merger, acquisition, or sale of assets, your information may be transferred as part of that transaction.</p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. International Data Transfers and Data Storage</h2>
              <p>
                We may store and process data using cloud infrastructure located in the United States, Europe, or other regions depending on our provider's architecture. Your data may be processed and stored on servers located outside your country of residence.
              </p>
              <p>
                When we transfer data internationally, we implement appropriate safeguards to protect your information in accordance with this policy and applicable data protection laws, including Standard Contractual Clauses (SCCs) where applicable.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Cookies and Similar Technologies</h2>
              <p>
                We use browser localStorage to temporarily store your shopping cart information. This information is deleted when you clear your browser storage or when you use the "Clear cart" function in our application.
              </p>
              <p>
                We do not use third-party tracking cookies or invasive tracking technologies.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">9. Your Rights (GDPR/LGPD)</h2>
              <p>Depending on your location, you have the following rights regarding your personal information:</p>
              
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Right of Access:</strong> You can request a copy of the personal information we have about you.</li>
                <li><strong>Right of Rectification:</strong> You can request that we correct any incorrect or incomplete information.</li>
                <li><strong>Right of Erasure:</strong> You can request that we delete your personal information under certain circumstances.</li>
                <li><strong>Right to Data Portability:</strong> You can request that we transfer your information to another service provider.</li>
                <li><strong>Right to Object:</strong> You can object to the processing of your personal information under certain circumstances.</li>
                <li><strong>Right to Restrict Processing:</strong> You can request that we restrict the processing of your information.</li>
                <li><strong>Right to Withdraw Consent:</strong> You can withdraw your consent at any time when processing is based on consent.</li>
                <li><strong>Right to Lodge a Complaint:</strong> You have the right to lodge a complaint with the data protection authority in your jurisdiction.</li>
              </ul>
              
              <p className="mt-4">
                To exercise any of these rights, please contact us using the information provided in the Contact section. We will respond to privacy-related requests within 30 days or less, as required by applicable law.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">10. Data Security</h2>
              <p>
                We implement reasonable technical and organizational security measures to protect your personal information against unauthorized access, alteration, disclosure, or destruction. These measures include:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Encryption of data in transit</li>
                <li>Secure storage of information</li>
                <li>Restricted access to personal information</li>
                <li>Regular monitoring of our security systems</li>
              </ul>
              <p className="mt-4">
                However, no method of transmission over the Internet or electronic storage is 100% secure. Although we strive to protect your information, we cannot guarantee its absolute security.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">10.1. Data Breach Notification</h2>
              <p>
                In the event of a data breach that may compromise personal information, we will:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Notify affected users without undue delay, and in any event within 72 hours where feasible</li>
                <li>Notify relevant data protection authorities as required by applicable regulations (e.g., within 72 hours under GDPR)</li>
                <li>Provide clear information about the nature of the breach, likely consequences, and measures taken to address it</li>
                <li>Take immediate steps to contain and remediate the breach</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">11. Children's Privacy</h2>
              <p>
                Our services are not directed to children under 18 years of age. We do not knowingly collect personal information from children. If we discover that we have collected information from a child without parental consent, we will take steps to delete that information from our systems.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">12. Changes to this Privacy Policy</h2>
              <p>
                We may update this Privacy Policy occasionally to reflect changes in our practices or for legal, operational, or regulatory reasons. We will notify you of any material changes by posting the new policy on this page and updating the "Last updated" date.
              </p>
              <p>
                We recommend that you review this policy periodically to stay informed about how we protect your information.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">13. Contact</h2>
              <p>
                If you have questions, concerns, or requests related to this Privacy Policy or the processing of your personal information, you can contact us through:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56957857558</li>
                <li><strong>Email:</strong> privacy@micartapro.com</li>
              </ul>
              <p className="mt-4">
                We will respond to your request within a reasonable time and in accordance with applicable data protection laws.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">14. Jurisdiction and Applicable Law</h2>
              <p>
                This Privacy Policy is governed by applicable data protection laws in your jurisdiction, including but not limited to the European Union's General Data Protection Regulation (GDPR), Brazil's General Data Protection Law (LGPD), and other applicable data protection laws.
              </p>
            </section>
          </div>

          <div className="mt-12 pt-8 border-t border-gray-200 flex flex-col sm:flex-row justify-between items-center gap-4">
            <Link 
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              ← Back to Home
            </Link>
            <Link 
              href="/terms"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              View Terms and Conditions →
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
