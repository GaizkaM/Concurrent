-- random_int_generator.adb
with Ada.Numerics.Discrete_Random;

package body Random_Int_Generator is

   subtype Int_Range is Integer range 0 .. 9;

   -- Instantiate the Discrete_Random package for the specified range
   package Random_Ints is new Ada.Numerics.Discrete_Random(Int_Range);
   use Random_Ints;

   -- Declare a generator object
   Gen : Generator;

   function Get_Random_Integer return Integer is
   begin
      return Random(Gen);
   end Get_Random_Integer;

begin
   -- Initialize the random number generator once when the package is elaborated
   Reset(Gen);

   -- Implementation of the Get_Random_Integer function

end Random_Int_Generator;
